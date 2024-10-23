!include "MUI2.nsh"

Name "Gravishken"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

; Define the name of the installer
OutFile "build\GravishkenSetup.exe"

; Set the default installation directory
InstallDir "$PROGRAMFILES\Gravishken"

; Default section start
Section "MainSection" SEC01
    ; Set output path to installation directory
    SetOutPath $INSTDIR

    ; Include your executable and DLL files
    File "build\gravishken.exe"
    File "build\urita.dll"
    File "build\WebView2Loader.dll"
    File "build\.env"

    ; Write the uninstaller executable
    WriteUninstaller "$INSTDIR\Uninstall.exe"

    ; CreateShortCut "<shortcut_path>" "<target_executable>" "<icon_path>" "<description>" "<working_directory>"
    ; Create a shortcut in the Start Menu
    CreateShortCut "$SMSTARTUP\Gravishken.lnk" "$INSTDIR\gravishken.exe"

    ; Optionally, create a shortcut on the Desktop
    CreateShortCut "$DESKTOP\Gravishken.lnk" "$INSTDIR\gravishken.exe"
SectionEnd

; Create an uninstaller
Section "Uninstall"
    Delete "$INSTDIR\gravishken.exe"
    Delete "$INSTDIR\urita.dll"
    Delete "$INSTDIR\WebView2Loader.dll"
    Delete "$INSTDIR\.env"
    Delete "$INSTDIR\Uninstall.exe"
    Delete "$SMSTARTUP\Gravishken.lnk"
    Delete "$DESKTOP\Gravishken.lnk"
    RMDir "$INSTDIR" ; Remove the directory if empty
SectionEnd
