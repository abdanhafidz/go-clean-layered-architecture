#Requires -Version 5.1
<#
.SYNOPSIS
    Automatic Dependency Injection Generator for Go Repositories

.DESCRIPTION
    Scans ./repositories/ directory, discovers all repository constructors,
    and generates provider/repositories_provider.go with full DI wiring.
    Repositories typically depend on database connection from ConfigProvider.

.EXAMPLE
    .\repository_injector.ps1
    
.NOTES
    - Works with PowerShell 5.1+ and PowerShell 7+
    - No external dependencies required
    - Supports multi-line constructor signatures
    - Repositories depend on database instance from ConfigProvider
#>

[CmdletBinding()]
param()

# Configuration
$RepositoriesDir = "./repositories"
$OutputFile = "provider/repositories_provider.go"
$ModulePath = "abdanhafidz.com/go-clean-layered-architecture/repositories"

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
class RepositoryInfo {
    [string]$ConstructorName    # NewAccountRepository
    [string]$Domain             # AccountRepository
    [string]$VarName            # accountRepository
    [System.Collections.Generic.List[Parameter]]$Parameters
    
    RepositoryInfo() {
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
    
    Write-ColorOutput "Scanning for repository constructors in $Directory..." "Cyan"
    
    if (-not (Test-Path $Directory)) {
        Write-ColorOutput "ERROR: Directory '$Directory' not found!" "Red"
        exit 1
    }
    
    $goFiles = Get-ChildItem -Path $Directory -Filter "*.go" -Recurse -File
    $repositories = [System.Collections.Generic.List[RepositoryInfo]]::new()
    
    foreach ($file in $goFiles) {
        $content = Get-Content $file.FullName -Raw
        
        # Match function signatures (support multi-line)
        # Pattern: func NewXxxRepository(...) XxxRepository
        $pattern = '(?ms)func\s+(New[a-zA-Z0-9]+Repository)\s*\(([^)]*)\)\s+([a-zA-Z0-9*_.]+Repository)'
        $matches = [regex]::Matches($content, $pattern)
        
        foreach ($match in $matches) {
            $constructorName = $match.Groups[1].Value
            $paramsStr = $match.Groups[2].Value
            $returnType = $match.Groups[3].Value
            
            # Extract domain name (XxxRepository)
            $domain = Normalize-TypeName $returnType
            $varName = Get-LowerCamelCase $domain
            
            $repo = [RepositoryInfo]::new()
            $repo.ConstructorName = $constructorName
            $repo.Domain = $domain
            $repo.VarName = $varName
            
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
                        $p.Name = "param$($repo.Parameters.Count)"
                        $p.RawType = $parts[0]
                    } else {
                        continue
                    }
                    
                    $p.NormalizedType = Normalize-TypeName $p.RawType
                    $repo.Parameters.Add($p)
                }
            }
            
            $repositories.Add($repo)
            Write-ColorOutput "  Found: $constructorName" "Green"
        }
    }
    
    if ($repositories.Count -eq 0) {
        Write-ColorOutput "No repository constructors found matching pattern 'NewXxxRepository'!" "Red"
        exit 1
    }
    
    Write-ColorOutput "`nTotal repositories discovered: $($repositories.Count)" "Blue"
    return $repositories
}

function Resolve-RepositoryArgument {
    param([Parameter]$Param)
    
    $type = $Param.NormalizedType
    $paramName = $Param.Name
    
    # DEPENDENCY RESOLUTION RULES FOR REPOSITORIES
    # ============================================
    
    # 1. Database connection patterns
    if ($type -match '^(DB|Database|Gorm|SqlDB|Connection)$' -or $paramName -match '^(db|database|conn|connection)$') {
        return "db"
    }
    
    # 2. *gorm.DB (most common in Go GORM projects)
    if ($Param.RawType -match 'gorm\.DB' -or $type -eq "DB") {
        return "db"
    }
    
    # 3. *sql.DB (standard library)
    if ($Param.RawType -match 'sql\.DB') {
        return "db"
    }
    
    # ADD MORE SPECIAL CASES HERE:
    # --------------------------------------------
    # Example: Redis connection
    # if ($type -eq "RedisClient" -or $paramName -match "redis") {
    #     return "redisClient"
    # }
    #
    # Example: MongoDB connection
    # if ($type -eq "MongoClient" -or $paramName -match "mongo") {
    #     return "mongoClient"
    # }
    #
    # Example: Cache
    # if ($type -eq "Cache" -or $paramName -match "cache") {
    #     return "cache"
    # }
    #
    # Example: Logger
    # if ($type -eq "Logger" -or $paramName -match "logger") {
    #     return "logger"
    # }
    # --------------------------------------------
    
    # 4. Fallback: unresolved type
    return "/* TODO: provide $($Param.RawType) */"
}

function Generate-ProviderCode {
    param([System.Collections.Generic.List[RepositoryInfo]]$Repositories)
    
    Write-ColorOutput "`nGenerating repositories provider code..." "Cyan"
    
    # Sort repositories alphabetically for consistent output
    $sortedRepos = $Repositories | Sort-Object -Property Domain
    
    $sb = [System.Text.StringBuilder]::new()
    [void]$sb.AppendLine("package provider")
    [void]$sb.AppendLine()
    [void]$sb.AppendLine("import `"$ModulePath`"")
    [void]$sb.AppendLine()
    
    # Interface
    [void]$sb.AppendLine("type RepositoriesProvider interface {")
    foreach ($repo in $sortedRepos) {
        $line = "`tProvide$($repo.Domain)() repositories.$($repo.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Struct
    [void]$sb.AppendLine("type repositoriesProvider struct {")
    foreach ($repo in $sortedRepos) {
        $line = "`t$($repo.VarName) repositories.$($repo.Domain)"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Constructor
    [void]$sb.AppendLine("func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {")
    [void]$sb.AppendLine("`tdbConfig := cfg.ProvideDatabaseConfig()")
    [void]$sb.AppendLine("`tdb := dbConfig.GetInstance()")
    [void]$sb.AppendLine()
    
    # Initialize repositories
    foreach ($repo in $sortedRepos) {
        $args = @()
        foreach ($param in $repo.Parameters) {
            $args += Resolve-RepositoryArgument $param
        }
        $argsStr = $args -join ", "
        $line = "`t$($repo.VarName) := repositories.$($repo.ConstructorName)($argsStr)"
        [void]$sb.AppendLine($line)
    }
    
    [void]$sb.AppendLine()
    [void]$sb.AppendLine("`treturn &repositoriesProvider{")
    foreach ($repo in $sortedRepos) {
        $line = "`t`t$($repo.VarName): $($repo.VarName),"
        [void]$sb.AppendLine($line)
    }
    [void]$sb.AppendLine("`t}")
    [void]$sb.AppendLine("}")
    [void]$sb.AppendLine()
    
    # Getter methods
    foreach ($repo in $sortedRepos) {
        [void]$sb.AppendLine("func (r *repositoriesProvider) Provide$($repo.Domain)() repositories.$($repo.Domain) {")
        [void]$sb.AppendLine("`treturn r.$($repo.VarName)")
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
    Write-ColorOutput "  Go Repository Provider Generator v1.0" "Blue"
    Write-ColorOutput "=========================================`n" "Blue"
    
    # Step 1: Parse all repository constructors
    $repositories = Parse-GoFiles -Directory $RepositoriesDir
    
    # Step 2: Generate provider code
    $code = Generate-ProviderCode -Repositories $repositories
    
    # Step 3: Write to file
    Write-ProviderFile -Code $code -OutputPath $OutputFile
    
    Write-ColorOutput "`nSUCCESS! Repositories provider generated successfully.`n" "Green"
    Write-ColorOutput "Next steps:" "Cyan"
    Write-ColorOutput "  1. Review $OutputFile" "White"
    Write-ColorOutput "  2. Fill any /* TODO: provide ... */ placeholders" "White"
    Write-ColorOutput "  3. Run: go build ./provider" "White"
    
} catch {
    Write-ColorOutput "`nERROR: $($_.Exception.Message)" "Red"
    Write-ColorOutput "Stack trace: $($_.ScriptStackTrace)" "Yellow"
    exit 1
}