#Requires -Version 5.1
<#
.SYNOPSIS
    Automatic Dependency Injection Generator for Go Middleware

.DESCRIPTION
    Scans ./middleware/ directory, discovers all middleware constructors, infers their dependencies,
    and generates provider/middleware_provider.go with full DI wiring.

.EXAMPLE
    .\middleware_injector.ps1
    
.NOTES
    - Works with PowerShell 5.1+ and PowerShell 7+
    - No external dependencies required
    - Supports multi-line constructor signatures
    - Middleware depend on services from ServicesProvider
#>

[CmdletBinding()]
param()

# Configuration
$MiddlewareDir = "./middleware"
$OutputFile = "provider/middleware_provider.go"
$ModulePath = "abdanhafidz.com/go-clean-layered-architecture/middleware"

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
class MiddlewareInfo {
    [string]$ConstructorName    # NewAuthenticationMiddleware
    [string]$Domain             # AuthenticationMiddleware
    [string]$VarName            # authenticationMiddleware
    [System.Collections.Generic.List[Parameter]]$Parameters
    
    MiddlewareInfo() {
        $this.Parameters = [System.Collections.Generic.List[Parameter]]::new()
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
    
    Write-ColorOutput "Scanning for middleware constructors in $Directory..." "Cyan"
    
    if (-not (Test-Path $Directory)) {
        Write-ColorOutput "ERROR: Directory '$Directory' not found!" "Red"
        exit 1
    }
    
    $goFiles = Get-ChildItem -Path $Directory -Filter "*.go" -Recurse -File
    $middlewares = [System.Collections.Generic.List[MiddlewareInfo]]::new()
    
    foreach ($file in $goFiles) {
        $content = Get-Content $file.FullName -Raw
        
        # Match function signatures (support multi-line)
        # Pattern: func NewXxxMiddleware(...) XxxMiddleware
        $pattern = '(?ms)func\s+(New[a-zA-Z0-9]+Middleware)\s*\(([^)]*)\)\s+([a-zA-Z0-9*_.]+Middleware)'
        $matches = [regex]::Matches($content, $pattern)
        
        foreach ($match in $matches) {
            $constructorName = $match.Groups[1].Value
            $paramsStr = $match.Groups[2].Value
            $returnType = $match.Groups[3].Value
            
            # Extract domain name (XxxMiddleware)
            $domain = Normalize-TypeName $returnType
            $varName = Get-LowerCamelCase $domain
            
            $middleware = [MiddlewareInfo]::new()
            $middleware.ConstructorName = $constructorName
            $middleware.Domain = $domain
            $middleware.VarName = $varName
            
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
                        $p.Name = "param$($middleware.Parameters.Count)"
                        $p.RawType = $parts[0]
                    } else {
                        continue
                    }
                    
                    $p.NormalizedType = Normalize-TypeName $p.RawType
                    $middleware.Parameters.Add($p)
                }
            }
            
            $middlewares.Add($middleware)
            Write-ColorOutput "  Found: $constructorName" "Green"
        }
    }
    
    if ($middlewares.Count -eq 0) {
        Write-ColorOutput "No middleware constructors found matching pattern 'NewXxxMiddleware'!" "Red"
        exit 1
    }
    
    Write-ColorOutput "`nTotal middleware discovered: $($middlewares.Count)" "Blue"
    return $middlewares
}

function Resolve-MiddlewareArgument {
    param([Parameter]$Param)
    
    $type = $Param.NormalizedType
    
    # DEPENDENCY RESOLUTION RULES
    # ============================================
    
    # 1. Service pattern: XxxxService -> servicesProvider.ProvideXxxxService()
    if ($type -match '^(.+)Service$') {
        $serviceName = $type
        return "servicesProvider.Provide${serviceName}()"
    }
    
    # ADD MORE SPECIAL CASES HERE:
    # --------------------------------------------
    # Example: Config dependency
    # if ($type -eq "Config") {
    #     return "configProvider.ProvideConfig()"
    # }
    #
    # Example: Logger
    # if ($type -eq "Logger") {
    #     return "loggerProvider.ProvideLogger()"
    # }
    #
    # Example: JWT Config
    # if ($type -eq "JWTConfig") {
    #     return "configProvider.ProvideJWTConfig()"
    # }
    #
    # Example: Database
    # if ($type -eq "DB" -or $type -eq "Database") {
    #     return "dbProvider.ProvideDatabase()"
    # }
    # --------------------------------------------
    
    # 2. Fallback: unresolved type
    return "/* TODO: provide $($Param.RawType) */"
}

function Generate-ProviderCode {
    param([System.Collections.Generic.List[MiddlewareInfo]]$Middlewares)
    
    Write-ColorOutput "`nGenerating middleware provider code..." "Cyan"
    
    # Sort middleware alphabetically for consistent output
    $sortedMiddlewares = $Middlewares | Sort-Object -Property Domain
    
    $sb = [System.Text.StringBuilder]::new()
    [void]$sb.AppendLine("package provider")
    [void]$sb.AppendLine()
    [void]$sb.AppendLine("import `"$ModulePath`"")
    [void]$sb.AppendLine()
    
    # Interface
    [void]$sb.AppendLine("type MiddlewareProvider interface {")
    foreach ($mw in $sortedMiddlewares) {
        $line = "`tProvide$($mw.Domain)() middleware.$($mw.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Struct
    [void]$sb.AppendLine("type middlewareProvider struct {")
    foreach ($mw in $sortedMiddlewares) {
        $line = "`t$($mw.VarName) middleware.$($mw.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Constructor
    [void]$sb.AppendLine("func NewMiddlewareProvider(servicesProvider ServicesProvider) MiddlewareProvider {")
    
    # Initialize middleware
    foreach ($mw in $sortedMiddlewares) {
        $args = @()
        foreach ($param in $mw.Parameters) {
            $args += Resolve-MiddlewareArgument $param
        }
        $argsStr = $args -join ", "
        $line = "`t$($mw.VarName) := middleware.$($mw.ConstructorName)($argsStr)"
        [void]$sb.AppendLine($line)
    }
    
    [void]$sb.AppendLine("`treturn &middlewareProvider{")
    foreach ($mw in $sortedMiddlewares) {
        $line = "`t`t$($mw.VarName): $($mw.VarName),"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("`t}")
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Getter methods
    foreach ($mw in $sortedMiddlewares) {
        [void]$sb.AppendLine("func (p *middlewareProvider) Provide$($mw.Domain)() middleware.$($mw.Domain) {")
        [void]$sb.AppendLine("`treturn p.$($mw.VarName)")
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
    Write-ColorOutput "  Go Middleware Provider Generator v1.0" "Blue"
    Write-ColorOutput "=========================================`n" "Blue"
    
    # Step 1: Parse all middleware constructors
    $middlewares = Parse-GoFiles -Directory $MiddlewareDir
    
    # Step 2: Generate provider code
    $code = Generate-ProviderCode -Middlewares $middlewares
    
    # Step 3: Write to file
    Write-ProviderFile -Code $code -OutputPath $OutputFile
    
    Write-ColorOutput "`nSUCCESS! Middleware provider generated successfully.`n" "Green"
    Write-ColorOutput "Next steps:" "Cyan"
    Write-ColorOutput "  1. Review $OutputFile" "White"
    Write-ColorOutput "  2. Fill any /* TODO: provide ... */ placeholders" "White"
    Write-ColorOutput "  3. Run: go build ./provider" "White"
    
} catch {
    Write-ColorOutput "`nERROR: $($_.Exception.Message)" "Red"
    Write-ColorOutput "Stack trace: $($_.ScriptStackTrace)" "Yellow"
    exit 1
}