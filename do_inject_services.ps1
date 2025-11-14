#Requires -Version 5.1
<#
.SYNOPSIS
    Automatic Dependency Injection Generator for Go Services

.DESCRIPTION
    Scans ./services/ directory, discovers service constructors, resolves dependencies,
    performs topological sorting, and generates provider/services_provider.go with full DI wiring.

.EXAMPLE
    .\service_injector.ps1
    
.NOTES
    - Works with PowerShell 5.1+ and PowerShell 7+
    - No external dependencies required
    - Supports multi-line constructor signatures
    - Handles dependency cycles with clear error messages
#>

[CmdletBinding()]
param()

# Configuration
$ServicesDir = "./services"
$OutputFile = "provider/services_provider.go"
$ModulePath = "abdanhafidz.com/go-clean-layered-architecture/services"

# ANSI colors for better output (fallback to plain text if not supported)
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
class ServiceInfo {
    [string]$ConstructorName    # NewAccountService
    [string]$Domain             # AccountService
    [string]$VarName            # accountService
    [System.Collections.Generic.List[Parameter]]$Parameters
    [System.Collections.Generic.List[string]]$ServiceDependencies
    
    ServiceInfo() {
        $this.Parameters = [System.Collections.Generic.List[Parameter]]::new()
        $this.ServiceDependencies = [System.Collections.Generic.List[string]]::new()
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
    
    Write-ColorOutput "Scanning for service constructors in $Directory..." "Cyan"
    
    if (-not (Test-Path $Directory)) {
        Write-ColorOutput "ERROR: Directory '$Directory' not found!" "Red"
        exit 1
    }
    
    $goFiles = Get-ChildItem -Path $Directory -Filter "*.go" -Recurse -File
    $services = [System.Collections.Generic.List[ServiceInfo]]::new()
    
    foreach ($file in $goFiles) {
        $content = Get-Content $file.FullName -Raw
        
        # Match function signatures (support multi-line)
        # Pattern: func NewXxxService(...) XxxService
        $pattern = '(?ms)func\s+(New[a-zA-Z0-9]+Service)\s*\(([^)]*)\)\s+([a-zA-Z0-9*_.]+Service)'
        $matches = [regex]::Matches($content, $pattern)
        
        foreach ($match in $matches) {
            $constructorName = $match.Groups[1].Value
            $paramsStr = $match.Groups[2].Value
            $returnType = $match.Groups[3].Value
            
            # Extract domain name (XxxService)
            $domain = Normalize-TypeName $returnType
            $varName = Get-LowerCamelCase $domain
            
            $service = [ServiceInfo]::new()
            $service.ConstructorName = $constructorName
            $service.Domain = $domain
            $service.VarName = $varName
            
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
                        $p.Name = "param$($service.Parameters.Count)"
                        $p.RawType = $parts[0]
                    } else {
                        continue
                    }
                    
                    $p.NormalizedType = Normalize-TypeName $p.RawType
                    $service.Parameters.Add($p)
                    
                    # Track service dependencies
                    if ($p.NormalizedType -match 'Service$') {
                        $service.ServiceDependencies.Add($p.NormalizedType)
                    }
                }
            }
            
            $services.Add($service)
            Write-ColorOutput "  Found: $constructorName" "Green"
        }
    }
    
    if ($services.Count -eq 0) {
        Write-ColorOutput "No service constructors found matching pattern 'NewXxxService'!" "Red"
        exit 1
    }
    
    Write-ColorOutput "`nTotal services discovered: $($services.Count)" "Blue"
    return $services
}

function Get-TopologicalOrder {
    param([System.Collections.Generic.List[ServiceInfo]]$Services)
    
    Write-ColorOutput "`nBuilding dependency graph..." "Cyan"
    
    # Build adjacency list
    $graph = @{}
    $inDegree = @{}
    $domainToService = @{}
    
    foreach ($svc in $Services) {
        $graph[$svc.Domain] = [System.Collections.Generic.List[string]]::new()
        $inDegree[$svc.Domain] = 0
        $domainToService[$svc.Domain] = $svc
    }
    
    # Build edges
    foreach ($svc in $Services) {
        foreach ($dep in $svc.ServiceDependencies) {
            if ($graph.ContainsKey($dep)) {
                $graph[$dep].Add($svc.Domain)
                $inDegree[$svc.Domain]++
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
    if ($sorted.Count -ne $Services.Count) {
        $remaining = $inDegree.Keys | Where-Object { $inDegree[$_] -gt 0 }
        Write-ColorOutput "`nERROR: Circular dependency detected!" "Red"
        Write-ColorOutput "Services involved in cycle: $($remaining -join ', ')" "Yellow"
        exit 1
    }
    
    Write-ColorOutput "  Dependency graph validated (no cycles)" "Green"
    Write-ColorOutput "  Topological order: $($sorted -join ' -> ')" "Blue"
    
    # Return services in topological order
    return $sorted | ForEach-Object { $domainToService[$_] }
}

function Resolve-ConstructorArgument {
    param([Parameter]$Param)
    
    $type = $Param.NormalizedType
    
    # SPECIAL CASE MAPPINGS - Add more here as needed
    # ============================================
    
    # 1. JWT secret string
    if ($Param.RawType -eq "string" -and $Param.Name -match "secret|key") {
        return "configProvider.ProvideJWTConfig().GetSecretKey()"
    }
    
    # 2. Repository pattern: XxxxRepository -> repoProvider.ProvideXxxxRepository()
    if ($type -match '^(.+)Repository$') {
        $repoName = $matches[1]
        return "repoProvider.Provide${repoName}Repository()"
    }
    
    # 3. Service pattern: XxxxService -> use variable (will be constructed before this)
    if ($type -match 'Service$') {
        return (Get-LowerCamelCase $type)
    }
    
    # ADD MORE SPECIAL CASES HERE:
    # --------------------------------------------
    # Example: Mail config
    # if ($type -eq "MailConfig") {
    #     return "configProvider.ProvideMailConfig()"
    # }
    #
    # Example: Redis client
    # if ($type -eq "RedisClient") {
    #     return "configProvider.ProvideRedisClient()"
    # }
    # --------------------------------------------
    
    # 4. Fallback: unresolved type
    return "/* TODO: provide $($Param.RawType) */"
}

function Generate-ProviderCode {
    param([System.Collections.Generic.List[ServiceInfo]]$ServicesInOrder)
    
    Write-ColorOutput "`nGenerating provider code..." "Cyan"
    
    $sb = [System.Text.StringBuilder]::new()
    [void]$sb.AppendLine("package provider")
    [void]$sb.AppendLine()
    [void]$sb.AppendLine("import `"$ModulePath`"")
    [void]$sb.AppendLine()
    
    # Interface
    [void]$sb.AppendLine("type ServicesProvider interface {")
    foreach ($svc in $ServicesInOrder) {
        $line = "`tProvide$($svc.Domain)() services.$($svc.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Struct
    [void]$sb.AppendLine("type servicesProvider struct {")
    foreach ($svc in $ServicesInOrder) {
        $line = "`t$($svc.VarName) services.$($svc.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Constructor
    [void]$sb.AppendLine("func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {")
    
    # Initialize services in topological order
    foreach ($svc in $ServicesInOrder) {
        $args = @()
        foreach ($param in $svc.Parameters) {
            $args += Resolve-ConstructorArgument $param
        }
        $argsStr = $args -join ", "
        $line = "`t$($svc.VarName) := services.$($svc.ConstructorName)($argsStr)"
        [void]$sb.AppendLine($line)
    }
    
    [void]$sb.AppendLine()
    [void]$sb.AppendLine("`treturn &servicesProvider{")
    foreach ($svc in $ServicesInOrder) {
        $line = "`t`t$($svc.VarName): $($svc.VarName),"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("`t}")
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Getter methods
    foreach ($svc in $ServicesInOrder) {
        [void]$sb.AppendLine("func (s *servicesProvider) Provide$($svc.Domain)() services.$($svc.Domain) {")
        [void]$sb.AppendLine("`treturn s.$($svc.VarName)")
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
    Write-ColorOutput "  Go Service Dependency Injector v1.0" "Blue"
    Write-ColorOutput "=========================================`n" "Blue"
    
    # Step 1: Parse all service constructors
    $services = Parse-GoFiles -Directory $ServicesDir
    
    # Step 2: Perform topological sort
    $sortedServices = Get-TopologicalOrder -Services $services
    
    # Step 3: Generate provider code
    $code = Generate-ProviderCode -ServicesInOrder $sortedServices
    
    # Step 4: Write to file
    Write-ProviderFile -Code $code -OutputPath $OutputFile
    
    Write-ColorOutput "`nSUCCESS! Provider generated successfully.`n" "Green"
    Write-ColorOutput "Next steps:" "Cyan"
    Write-ColorOutput "  1. Review $OutputFile" "White"
    Write-ColorOutput "  2. Fill any /* TODO: provide ... */ placeholders" "White"
    Write-ColorOutput "  3. Run: go build ./provider" "White"
    
} catch {
    Write-ColorOutput "`nERROR: $($_.Exception.Message)" "Red"
    Write-ColorOutput "Stack trace: $($_.ScriptStackTrace)" "Yellow"
    exit 1
}