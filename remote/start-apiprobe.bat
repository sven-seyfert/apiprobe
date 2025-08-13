cls

@REM Guard logic
@REM Do not execute main logic within timespan 22:00 - 06:10.
@REM That's only an example on how to do such exception time.
for /F "tokens=1-3 delims=:" %%a in ("%time%") do (
    set hour=%%a
    set minute=%%b
)

@REM Remove whitespaces
set hour=%hour: =%
set minute=%minute: =%

@REM Calculation of hours and minutes in minutes
set /a currentTimeInMinutes=%hour% * 60 + %minute%

@REM Exception timespan 22:00 - 06:10
set /a startExceptionTime=22 * 60
set /a endExceptionTime=6 * 60 + 10

@REM Exit cases
if %currentTimeInMinutes% geq %startExceptionTime% (
    exit
)

if %currentTimeInMinutes% leq %endExceptionTime% (
    exit
)

@REM Main logic
cls

@REM Remove previously opened instance in case it
@REM was not quit correctly.
taskkill /F /FI "WINDOWTITLE eq cmd-apiprobe" /T

@REM Give this instance a name.
title cmd-apiprobe

@REM This is the path were the project should be located.
@REM Example path on the target maschine.
cd C:\Store\Repositories\GitHub\apiprobe

@REM Run tests for different environments (by tag)
@REM and exclude specific request by ID (for test environment).
call apiprobe.exe --name "Environment: TEST" --tags "env-test" --exclude "ff00fceb61"
call apiprobe.exe --name "Environment: PROD" --tags "env-prod"

@REM Increase robustness of finishing previous processing
ping localhost -n 10

exit
