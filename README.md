# 用法
```cmd
path-cleaner [-path  your_path_str ]  remove_1 remove_2 remove_3 ...  
```

## 用途
这个主要是用来在 windows 的命令行时，切换 jdk 时用来清除多途的JDK 路径的，在反复SET PATH 之后，PATH 会越来越长，这个是清理PATH 路径的。

## JDK 切换脚本示例

### cmd版本
```bat
@echo off
setlocal

:: ==============================
:: 1. 配置路径
:: ==============================
set "JAVA8_HOME=C:\Program Files\Eclipse Adoptium\jdk-8.0.432.6-hotspot"
set "JAVA8_BIN=%JAVA8_HOME%\bin"

set "JAVA21_HOME=C:\Program Files\Eclipse Adoptium\jdk-21.0.5.11-hotspot"
set "JAVA21_BIN=%JAVA21_HOME%\bin"

:: ==============================
:: 2. 参数判断
:: ==============================
:: 如果没有参数，跳转到信息展示
if "%~1"=="" goto INFO

:: 根据参数跳转
if "%~1"=="8" goto SET_JAVA8
if "%~1"=="1.8" goto SET_JAVA8
if "%~1"=="21" goto SET_JAVA21

:: 如果参数无法识别
echo.
echo [Error] Unknown version: %~1
goto INFO

:: ==============================
:: 3. 信息展示 (无参数时)
:: ==============================
:INFO
echo.
echo ============================================
echo           Java Switcher Info
echo ============================================
echo Current Java Version:
java -version 2>&1 | findstr "version"
echo.
echo Available Versions:
echo   [8]  : Eclipse Adoptium JDK 8
echo   [21] : Eclipse Adoptium JDK 21
echo.
echo Usage:
echo   switch-java       (Show this info)
echo   switch-java 8     (Switch to Java 8)
echo   switch-java 21    (Switch to Java 21)
echo ============================================
goto END

:: ==============================
:: 4. 设定目标变量
:: ==============================
:SET_JAVA8
set "NEW_JAVA_HOME=%JAVA8_HOME%"
set "NEW_JAVA_BIN=%JAVA8_BIN%"
set "TARGET_VER=8"
goto ACTIVATE

:SET_JAVA21
set "NEW_JAVA_HOME=%JAVA21_HOME%"
set "NEW_JAVA_BIN=%JAVA21_BIN%"
set "TARGET_VER=21"
goto ACTIVATE

:: ==============================
:: 5. 执行切换 (集成 path-cleaner)
:: ==============================
:ACTIVATE
:: 调用 path-cleaner 清理旧路径
for /f "delims=" %%i in ('path-cleaner "%JAVA8_BIN%" "%JAVA21_BIN%"') do set "CLEANED_PATH=%%i"

:: 防止 cleaner 异常导致 path 丢失
if "%CLEANED_PATH%"=="" (
    echo [Error] path-cleaner failed. Aborting.
    goto END
)

:: 应用环境变量 (逃逸到外部环境)
endlocal & set "JAVA_HOME=%NEW_JAVA_HOME%" & set "Path=%NEW_JAVA_BIN%;%CLEANED_PATH%"

echo.
echo [Success] Switched to Java %TARGET_VER%
java -version 2>&1 | findstr "version"

:END

```


### powershell 函数，写在`$Profile` 中的，类似 linux 的 `.bashrc`

```ps1
function switch-java {
    param (
        [string]$Version
    )

    # --- 1. 配置路径 ---
    $Java8Home  = "C:\Program Files\Eclipse Adoptium\jdk-8.0.432.6-hotspot"
    $Java21Home = "C:\Program Files\Eclipse Adoptium\jdk-21.0.5.11-hotspot"
    
    $Bin8  = "$Java8Home\bin"
    $Bin21 = "$Java21Home\bin"

    # --- 2. 无参数模式：显示信息 ---
    if ([string]::IsNullOrWhiteSpace($Version)) {
        Write-Host "`n============================================" -ForegroundColor Cyan
        Write-Host "          Java Switcher Info" -ForegroundColor Cyan
        Write-Host "============================================" -ForegroundColor Cyan
        Write-Host "Current Java:" -ForegroundColor Gray
        java -version
        Write-Host "`nAvailable Commands:" -ForegroundColor Gray
        Write-Host "  switch-java       " -ForegroundColor Yellow -NoNewline; Write-Host "(Show this info)"
        Write-Host "  switch-java 8     " -ForegroundColor Yellow -NoNewline; Write-Host "(Switch to Java 8)"
        Write-Host "  switch-java 21    " -ForegroundColor Yellow -NoNewline; Write-Host "(Switch to Java 21)"
        Write-Host "============================================" -ForegroundColor Cyan
        return
    }

    # --- 3. 确定目标 ---
    $TargetHome = ""
    $TargetBin  = ""

    if ($Version -eq "8" -or $Version -eq "1.8") {
        $TargetHome = $Java8Home
        $TargetBin  = $Bin8
    } elseif ($Version -eq "21") {
        $TargetHome = $Java21Home
        $TargetBin  = $Bin21
    } else {
        Write-Error "Unknown version: $Version. Please use '8' or '21'."
        return
    }

    # --- 4. 执行切换 (集成 path-cleaner) ---
    try {
        # 调用 path-cleaner
        $CleanedPath = path-cleaner "$Bin8" "$Bin21"
        
        if ([string]::IsNullOrWhiteSpace($CleanedPath)) {
            Write-Error "path-cleaner returned empty string. Aborting."
            return
        }

        # 设置环境变量
        $env:JAVA_HOME = $TargetHome
        $env:Path = "$TargetBin;$CleanedPath"

        Write-Host "`n[Success] Switched to Java $Version" -ForegroundColor Green
        java -version
    }
    catch {
        Write-Error "Failed to run path-cleaner. Make sure it is in your Path."
    }
}
```
