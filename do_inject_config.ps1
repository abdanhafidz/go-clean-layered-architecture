#Requires -Version 5.1
<#
.SYNOPSIS
    Automatic Dependency Injection Generator for Go Configuration

.DESCRIPTION
    Scans ./config/ directory, discovers all config constructors, infers their dependencies,
    and generates provider/config_provider.go with full DI wiring. Special handling for
    chained dependencies where configs depend on other configs (e.g., DatabaseConfig depends on EnvConfig).

.EXAMPLE
    .\config_injector.ps1
    
.NOTES
    - Works with PowerShell 5.1+ and PowerShell 7+
    - No external dependencies required
    - Supports multi-line constructor signatures
    - Handles config-to-config dependencies with topological sorting
    - Special handling for method calls like envConfig.GetDatabaseHost()
#>

[CmdletBinding()]
param()

# Configuration
$ConfigDir = "./config"
$OutputFile = "provider/config_provider.go"
$ModulePath = "abdanhafidz.com/go-clean-layered-architecture/config"

# ANSI colors for better output
$script:UseColors = $Host.UI.SupportsVirtualTerminal
function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    if ($script:UseColors) {
        $colors = @{
            "Green" = "`e[32m"; "Yellow" = "`e[33m"; "Red" = "`e[31m"
            "Cyan" = "`e[36m"; "Blue" = "`e[34m"; "Reset" = "`e[0m"
        }
        Write-Host "$($colors[$Color])$Message$($colors['Reset'])"
    } else {
        Write-Host $Message
    }
}

# Data structures
class ConfigInfo {
    [string]$ConstructorName    # NewDatabaseConfig
    [string]$Domain             # DatabaseConfig
    [string]$VarName            # databaseConfig
    [System.Collections.Generic.List[Parameter]]$Parameters
    [System.Collections.Generic.List[string]]$ConfigDependencies
    
    ConfigInfo() {
        $this.Parameters = [System.Collections.Generic.List[Parameter]]::new()
        $this.ConfigDependencies = [System.Collections.Generic.List[string]]::new()
    }
}

class Parameter {
    [string]$Name
    [string]$RawType
    [string]$NormalizedType
}

function Get-LowerCamelCase {
    param([string]$Text)
    if ($Text.Length -eq 0) { return $Text }
    return $Text.Substring(0, 1).ToLower() + $Text.Substring(1)
}

function Normalize-TypeName {
    param([string]$TypeStr)
    
    # Remove leading pointer
    $cleaned = $TypeStr -replace '^\*+', ''
    
    # Remove package prefix (everything before last dot)
    if ($cleaned -match '\.([^.]+)$') {
        $cleaned = $matches[1]
    }
    
    return $cleaned.Trim()
}

function Parse-GoFiles {
    param([string]$Directory)
    
    Write-ColorOutput "Scanning for config constructors in $Directory..." "Cyan"
    
    if (-not (Test-Path $Directory)) {
        Write-ColorOutput "ERROR: Directory '$Directory' not found!" "Red"
        exit 1
    }
    
    $goFiles = Get-ChildItem -Path $Directory -Filter "*.go" -Recurse -File
    $configs = [System.Collections.Generic.List[ConfigInfo]]::new()
    
    foreach ($file in $goFiles) {
        $content = Get-Content $file.FullName -Raw
        
        # Match function signatures (support multi-line)
        # Pattern: func NewXxxConfig(...) XxxConfig
        $pattern = '(?ms)func\s+(New[a-zA-Z0-9]+Config)\s*\(([^)]*)\)\s+([a-zA-Z0-9*_.]+Config)'
        $matches = [regex]::Matches($content, $pattern)
        
        foreach ($match in $matches) {
            $constructorName = $match.Groups[1].Value
            $paramsStr = $match.Groups[2].Value
            $returnType = $match.Groups[3].Value
            
            # Extract domain name (XxxConfig)
            $domain = Normalize-TypeName $returnType
            $varName = Get-LowerCamelCase $domain
            
            $config = [ConfigInfo]::new()
            $config.ConstructorName = $constructorName
            $config.Domain = $domain
            $config.VarName = $varName
            
            # Parse parameters
            if ($paramsStr.Trim() -ne "") {
                # Split by comma, but be careful with nested types
                $paramList = $paramsStr -split ',\s*(?![^<>]*>)'
                
                foreach ($param in $paramList) {
                    $param = $param.Trim()
                    if ($param -eq "") { continue }
                    
                    # Split into name and type
                    $parts = $param -split '\s+', 2
                    
                    $p = [Parameter]::new()
                    if ($parts.Count -eq 2) {
                        $p.Name = $parts[0]
                        $p.RawType = $parts[1]
                    } elseif ($parts.Count -eq 1) {
                        # Anonymous parameter - synthesize name
                        $p.Name = "param$($config.Parameters.Count)"
                        $p.RawType = $parts[0]
                    } else {
                        continue
                    }
                    
                    $p.NormalizedType = Normalize-TypeName $p.RawType
                    $config.Parameters.Add($p)
                    
                    # Track config dependencies (configs that depend on other configs)
                    if ($p.NormalizedType -match 'Config$') {
                        $config.ConfigDependencies.Add($p.NormalizedType)
                    }
                }
            }
            
            $configs.Add($config)
            Write-ColorOutput "  Found: $constructorName" "Green"
        }
    }
    
    if ($configs.Count -eq 0) {
        Write-ColorOutput "No config constructors found matching pattern 'NewXxxConfig'!" "Red"
        exit 1
    }
    
    Write-ColorOutput "`nTotal configs discovered: $($configs.Count)" "Blue"
    return $configs
}

function Get-TopologicalOrder {
    param([System.Collections.Generic.List[ConfigInfo]]$Configs)
    
    Write-ColorOutput "`nBuilding dependency graph..." "Cyan"
    
    # Build adjacency list
    $graph = @{}
    $inDegree = @{}
    $domainToConfig = @{}
    
    foreach ($cfg in $Configs) {
        $graph[$cfg.Domain] = [System.Collections.Generic.List[string]]::new()
        $inDegree[$cfg.Domain] = 0
        $domainToConfig[$cfg.Domain] = $cfg
    }
    
    # Build edges (dependencies)
    foreach ($cfg in $Configs) {
        foreach ($dep in $cfg.ConfigDependencies) {
            if ($graph.ContainsKey($dep)) {
                $graph[$dep].Add($cfg.Domain)
                $inDegree[$cfg.Domain]++
            }
        }
    }
    
    # Kahn's algorithm for topological sort
    $queue = [System.Collections.Generic.Queue[string]]::new()
    foreach ($domain in $inDegree.Keys) {
        if ($inDegree[$domain] -eq 0) {
            $queue.Enqueue($domain)
        }
    }
    
    $sorted = [System.Collections.Generic.List[string]]::new()
    
    while ($queue.Count -gt 0) {
        $current = $queue.Dequeue()
        $sorted.Add($current)
        
        foreach ($neighbor in $graph[$current]) {
            $inDegree[$neighbor]--
            if ($inDegree[$neighbor] -eq 0) {
                $queue.Enqueue($neighbor)
            }
        }
    }
    
    # Check for cycles
    if ($sorted.Count -ne $Configs.Count) {
        $remaining = $inDegree.Keys | Where-Object { $inDegree[$_] -gt 0 }
        Write-ColorOutput "`nERROR: Circular dependency detected in configs!" "Red"
        Write-ColorOutput "Configs involved in cycle: $($remaining -join ', ')" "Yellow"
        exit 1
    }
    
    Write-ColorOutput "  Dependency graph validated (no cycles)" "Green"
    Write-ColorOutput "  Topological order: $($sorted -join ' -> ')" "Blue"
    
    # Return configs in topological order
    return $sorted | ForEach-Object { $domainToConfig[$_] }
}

function Resolve-ConfigArgument {
    param(
        [Parameter]$Param,
        [string]$DependentConfigVar
    )
    
    $type = $Param.NormalizedType
    $paramName = $Param.Name
    
    # SPECIAL CASE MAPPINGS FOR CONFIG
    # ============================================
    
    # 1. Config pattern: XxxConfig -> already instantiated config variable
    if ($type -match '^(.+)Config$') {
        $configVarName = Get-LowerCamelCase $type
        return $configVarName
    }
    
    # 2. String parameters - try to infer from parameter name and match with config getter methods
    if ($type -eq "string") {
        # Common patterns for EnvConfig getters
        $getterMappings = @{
            "host" = "GetDatabaseHost()"
            "databaseHost" = "GetDatabaseHost()"
            "user" = "GetDatabaseUser()"
            "databaseUser" = "GetDatabaseUser()"
            "password" = "GetDatabasePassword()"
            "databasePassword" = "GetDatabasePassword()"
            "name" = "GetDatabaseName()"
            "databaseName" = "GetDatabaseName()"
            "dbName" = "GetDatabaseName()"
            "port" = "GetDatabasePort()"
            "databasePort" = "GetDatabasePort()"
            "salt" = "GetSalt()"
            "secret" = "GetSecretKey()"
            "secretKey" = "GetSecretKey()"
            "jwtSecret" = "GetSecretKey()"
            "apiKey" = "GetAPIKey()"
            "timezone" = "GetTimezone()"
        }
        
        # Try to find matching getter
        foreach ($key in $getterMappings.Keys) {
            if ($paramName -like "*$key*") {
                return "envConfig.$($getterMappings[$key])"
            }
        }
        
        # If dependent on a config, try to construct getter name from param name
        if ($DependentConfigVar) {
            # Convert paramName to PascalCase for getter
            $getterName = (Get-Culture).TextInfo.ToTitleCase($paramName)
            $getterName = $getterName -replace '\s', ''
            return "${DependentConfigVar}.Get${getterName}()"
        }
    }
    
    # 3. Int/port parameters
    if ($type -eq "int" -or $type -eq "int32" -or $type -eq "int64") {
        if ($paramName -match "port") {
            return "envConfig.GetDatabasePort()"
        }
    }
    
    # ADD MORE SPECIAL CASES HERE:
    # --------------------------------------------
    # Example: Redis config
    # if ($paramName -match "redis") {
    #     return "envConfig.GetRedisURL()"
    # }
    #
    # Example: Mail config
    # if ($paramName -match "smtp") {
    #     return "envConfig.GetSMTPHost()"
    # }
    # --------------------------------------------
    
    # 4. Hardcoded constants (timezone example)
    if ($type -eq "string" -and $paramName -match "timezone|location") {
        return "`"Asia/Jakarta`""
    }
    
    # 5. Fallback: unresolved type
    return "/* TODO: provide $($Param.RawType) for $paramName */"
}

function Generate-ProviderCode {
    param([System.Collections.Generic.List[ConfigInfo]]$ConfigsInOrder)
    
    Write-ColorOutput "`nGenerating config provider code..." "Cyan"
    
    $sb = [System.Text.StringBuilder]::new()
    [void]$sb.AppendLine("package provider")
    [void]$sb.AppendLine()
    [void]$sb.AppendLine("import `"$ModulePath`"")
    [void]$sb.AppendLine()
    
    # Interface
    [void]$sb.AppendLine("type ConfigProvider interface {")
    foreach ($cfg in $ConfigsInOrder) {
        $line = "`tProvide$($cfg.Domain)() config.$($cfg.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Struct
    [void]$sb.AppendLine("type configProvider struct {")
    foreach ($cfg in $ConfigsInOrder) {
        $line = "`t$($cfg.VarName) config.$($cfg.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Constructor
    [void]$sb.AppendLine("func NewConfigProvider() ConfigProvider {")
    
    # Initialize configs in topological order
    foreach ($cfg in $ConfigsInOrder) {
        # Check if this config depends on another config
        $dependentConfigVar = $null
        if ($cfg.ConfigDependencies.Count -gt 0) {
            $dependentConfigVar = Get-LowerCamelCase $cfg.ConfigDependencies[0]
        }
        
        $args = @()
        foreach ($param in $cfg.Parameters) {
            $args += Resolve-ConfigArgument -Param $param -DependentConfigVar $dependentConfigVar
        }
        $argsStr = $args -join ", "
        $line = "`t$($cfg.VarName) := config.$($cfg.ConstructorName)($argsStr)"
        [void]$sb.AppendLine($line)
    }
    
    [void]$sb.AppendLine("`treturn &configProvider{")
    foreach ($cfg in $ConfigsInOrder) {
        $line = "`t`t$($cfg.VarName): $($cfg.VarName),"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("`t}")
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Getter methods
    foreach ($cfg in $ConfigsInOrder) {
        [void]$sb.AppendLine("func (c *configProvider) Provide$($cfg.Domain)() config.$($cfg.Domain) {")
        [void]$sb.AppendLine("`treturn c.$($cfg.VarName)")
        [void]$sb.AppendLine("}")
        [void]$sb.AppendLine()
    }
    
    return $sb.ToString()
}

function Write-ProviderFile {
    param([string]$Code, [string]$OutputPath)
    
    Write-ColorOutput "Writing to $OutputPath..." "Cyan"
    
    # Ensure directory exists
    $dir = Split-Path $OutputPath -Parent
    if ($dir -and -not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
    }
    
    # Write file as UTF-8 without BOM
    $utf8NoBom = [System.Text.UTF8Encoding]::new($false)
    [System.IO.File]::WriteAllText($OutputPath, $Code, $utf8NoBom)
    
    Write-ColorOutput "  Successfully generated $OutputPath" "Green"
}

# ============================================
# MAIN EXECUTION
# ============================================

try {
    Write-ColorOutput "`n=========================================" "Blue"
    Write-ColorOutput "  Go Config Provider Generator v1.0" "Blue"
    Write-ColorOutput "=========================================`n" "Blue"
    
    # Step 1: Parse all config constructors
    $configs = Parse-GoFiles -Directory $ConfigDir
    
    # Step 2: Perform topological sort (configs can depend on other configs)
    $sortedConfigs = Get-TopologicalOrder -Configs $configs
    
    # Step 3: Generate provider code
    $code = Generate-ProviderCode -ConfigsInOrder $sortedConfigs
    
    # Step 4: Write to file
    Write-ProviderFile -Code $code -OutputPath $OutputFile
    
    Write-ColorOutput "`nSUCCESS! Config provider generated successfully.`n" "Green"
    Write-ColorOutput "Next steps:" "Cyan"
    Write-ColorOutput "  1. Review $OutputFile" "White"
    Write-ColorOutput "  2. Fill any /* TODO: provide ... */ placeholders" "White"
    Write-ColorOutput "  3. Run: go build ./provider" "White"
    
} catch {
    Write-ColorOutput "`nERROR: $($_.Exception.Message)" "Red"
    Write-ColorOutput "Stack trace: $($_.ScriptStackTrace)" "Yellow"
    exit 1
}