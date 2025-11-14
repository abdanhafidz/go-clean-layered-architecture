#Requires -Version 5.1
<#
.SYNOPSIS
    Automatic Dependency Injection Generator for Go Controllers

.DESCRIPTION
    Scans ./controllers/ directory, discovers all controller constructors, infers their dependencies,
    and generates provider/controller_provider.go with full DI wiring.

.EXAMPLE
    .\controller_injector.ps1
    
.NOTES
    - Works with PowerShell 5.1+ and PowerShell 7+
    - No external dependencies required
    - Supports multi-line constructor signatures
    - Controllers depend on services from ServicesProvider
#>

[CmdletBinding()]
param()

# Configuration
$ControllersDir = "./controllers"
$OutputFile = "provider/controller_provider.go"
$ModulePath = "abdanhafidz.com/go-clean-layered-architecture/controllers"

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
class ControllerInfo {
    [string]$ConstructorName    # NewAccountController
    [string]$Domain             # AccountController
    [string]$VarName            # accountController
    [System.Collections.Generic.List[Parameter]]$Parameters
    
    ControllerInfo() {
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
    
    Write-ColorOutput "Scanning for controller constructors in $Directory..." "Cyan"
    
    if (-not (Test-Path $Directory)) {
        Write-ColorOutput "ERROR: Directory '$Directory' not found!" "Red"
        exit 1
    }
    
    $goFiles = Get-ChildItem -Path $Directory -Filter "*.go" -Recurse -File
    $controllers = [System.Collections.Generic.List[ControllerInfo]]::new()
    
    foreach ($file in $goFiles) {
        $content = Get-Content $file.FullName -Raw
        
        # Match function signatures (support multi-line)
        # Pattern: func NewXxxController(...) XxxController
        $pattern = '(?ms)func\s+(New[a-zA-Z0-9]+Controller)\s*\(([^)]*)\)\s+([a-zA-Z0-9*_.]+Controller)'
        $matches = [regex]::Matches($content, $pattern)
        
        foreach ($match in $matches) {
            $constructorName = $match.Groups[1].Value
            $paramsStr = $match.Groups[2].Value
            $returnType = $match.Groups[3].Value
            
            # Extract domain name (XxxController)
            $domain = Normalize-TypeName $returnType
            $varName = Get-LowerCamelCase $domain
            
            $controller = [ControllerInfo]::new()
            $controller.ConstructorName = $constructorName
            $controller.Domain = $domain
            $controller.VarName = $varName
            
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
                        $p.Name = "param$($controller.Parameters.Count)"
                        $p.RawType = $parts[0]
                    } else {
                        continue
                    }
                    
                    $p.NormalizedType = Normalize-TypeName $p.RawType
                    $controller.Parameters.Add($p)
                }
            }
            
            $controllers.Add($controller)
            Write-ColorOutput "  Found: $constructorName" "Green"
        }
    }
    
    if ($controllers.Count -eq 0) {
        Write-ColorOutput "No controller constructors found matching pattern 'NewXxxController'!" "Red"
        exit 1
    }
    
    Write-ColorOutput "`nTotal controllers discovered: $($controllers.Count)" "Blue"
    return $controllers
}

function Resolve-ControllerArgument {
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
    # Example: Validator
    # if ($type -eq "Validator") {
    #     return "validatorProvider.ProvideValidator()"
    # }
    # --------------------------------------------
    
    # 2. Fallback: unresolved type
    return "/* TODO: provide $($Param.RawType) */"
}

function Generate-ProviderCode {
    param([System.Collections.Generic.List[ControllerInfo]]$Controllers)
    
    Write-ColorOutput "`nGenerating controller provider code..." "Cyan"
    
    # Sort controllers alphabetically for consistent output
    $sortedControllers = $Controllers | Sort-Object -Property Domain
    
    $sb = [System.Text.StringBuilder]::new()
    [void]$sb.AppendLine("package provider")
    [void]$sb.AppendLine()
    [void]$sb.AppendLine("import `"$ModulePath`"")
    [void]$sb.AppendLine()
    
    # Interface
    [void]$sb.AppendLine("type ControllerProvider interface {")
    foreach ($ctrl in $sortedControllers) {
        $line = "`tProvide$($ctrl.Domain)() controllers.$($ctrl.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Struct
    [void]$sb.AppendLine("type controllerProvider struct {")
    foreach ($ctrl in $sortedControllers) {
        $line = "`t$($ctrl.VarName) controllers.$($ctrl.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Constructor
    [void]$sb.AppendLine("func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {")
    [void]$sb.AppendLine()
    
    # Initialize controllers
    foreach ($ctrl in $sortedControllers) {
        $args = @()
        foreach ($param in $ctrl.Parameters) {
            $args += Resolve-ControllerArgument $param
        }
        $argsStr = $args -join ", "
        $line = "`t$($ctrl.VarName) := controllers.$($ctrl.ConstructorName)($argsStr)"
        [void]$sb.AppendLine($line)
    }
    
    [void]$sb.AppendLine("`treturn &controllerProvider{")
    foreach ($ctrl in $sortedControllers) {
        $line = "`t`t$($ctrl.VarName): $($ctrl.VarName),"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("`t}")
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Getter methods
    [void]$sb.AppendLine("// --- Getter Methods ---")
    [void]$sb.AppendLine()
    foreach ($ctrl in $sortedControllers) {
        [void]$sb.AppendLine("func (c *controllerProvider) Provide$($ctrl.Domain)() controllers.$($ctrl.Domain) {")
        [void]$sb.AppendLine("`treturn c.$($ctrl.VarName)")
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
    Write-ColorOutput "  Go Controller Provider Generator v1.0" "Blue"
    Write-ColorOutput "=========================================`n" "Blue"
    
    # Step 1: Parse all controller constructors
    $controllers = Parse-GoFiles -Directory $ControllersDir
    
    # Step 2: Generate provider code
    $code = Generate-ProviderCode -Controllers $controllers
    
    # Step 3: Write to file
    Write-ProviderFile -Code $code -OutputPath $OutputFile
    
    Write-ColorOutput "`nSUCCESS! Controller provider generated successfully.`n" "Green"
    Write-ColorOutput "Next steps:" "Cyan"
    Write-ColorOutput "  1. Review $OutputFile" "White"
    Write-ColorOutput "  2. Fill any /* TODO: provide ... */ placeholders" "White"
    Write-ColorOutput "  3. Run: go build ./provider" "White"
    
} catch {
    Write-ColorOutput "`nERROR: $($_.Exception.Message)" "Red"
    Write-ColorOutput "Stack trace: $($_.ScriptStackTrace)" "Yellow"
    exit 1
}