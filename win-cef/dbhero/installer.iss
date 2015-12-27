; using https://github.com/stfx/innodependencyinstaller for net 4.5 and vs 2013 x86 redist installations
; more info:
;  http://stackoverflow.com/questions/16054969/inno-setup-how-can-i-put-a-message-on-the-welcome-page
;  http://stackoverflow.com/questions/27587827/how-to-add-new-text-in-welcome-page-in-blank-space-by-changing-inno-setup-script
;  http://blog.elangroup-software.com/2012/08/inno-setup-part-8-creating-custom.html
;  http://stackoverflow.com/questions/4311995/inno-setup-custom-page
;  http://stackoverflow.com/questions/22183811/how-to-skip-all-the-wizard-pages-and-go-directly-to-the-installation-process
;  http://stackoverflow.com/questions/13921535/skipping-custom-pages-based-on-optional-components-in-inno-setup
;  http://www.jrsoftware.org/ishelp/
;
; TODO: fully custom page with just "Install" button
; TODO: always create desktop icon

#define MyAppName "dbHero"
;MyAppVersion is set from cmd-line
;#define MyAppVersion "9.9"
#define MyAppPublisher "Database Experts"
#define MyAppURL "http://www.dbheroapp.com"
#define MyAppExeName "dbHero.exe"

[Setup]
; NOTE: The value of AppId uniquely identifies this application.
; Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={{784A5D3A-FB9F-4E09-809C-40F639F408D6}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
; name shown in control panel's list of installed apps
; some apps show version there as well but it's redundant
; as it's shown anyway in another column
AppVerName={#MyAppName}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={localappdata}\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
OutputDir=bin\Release
OutputBaseFilename=dbHero-setup-inno
SetupIconFile=icon.ico
Compression=lzma
SolidCompression=yes
; Vista SP 1, http://www.jrsoftware.org/ishelp/index.php?topic=winvernotes
; http://www.gaijin.at/en/lstwinver.php
MinVersion=6.0.6002
DisableDirPage=yes
;DisableWelcomePage=yes
;DisableReadyPage=yes
UninstallDisplayIcon={app}\{#MyAppExeName}

[Languages]
Name: "en"; MessagesFile: "compiler:Default.isl"
Name: "de"; MessagesFile: "compiler:Languages\German.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"

[Dirs]
Name: "{app}\locales"

[Files]
; NOTE: Don't use "Flags: ignoreversion" on any shared system files
Source: "bin\Release\dbHero.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\dbHero.exe.config"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\Yepi.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "dbherohelper.exe"; DestDir: "{app}"; Flags: ignoreversion

; https://github.com/cefsharp/CefSharp/wiki/Output-files-description-table-%28Redistribution%29
Source: "bin\Release\CefSharp.BrowserSubprocess.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\CefSharp.BrowserSubprocess.Core.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\CefSharp.Core.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\CefSharp.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\CefSharp.WinForms.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\d3dcompiler_47.dll"; DestDir: "{app}"; Flags: ignoreversion
; only needed for html5 audio/video support
; Source: "bin\Release\ffmpegsumo.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\libcef.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\libEGL.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\libGLESv2.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\icudtl.dat"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\cef.pak"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\cef_extensions.pak"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\cef_100_percent.pak"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\cef_200_percent.pak"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\natives_blob.bin"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\locales\*"; DestDir: "{app}\locales"; Flags: ignoreversion

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{commondesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent

; based on https://github.com/stfx/innodependencyinstaller/blob/master/setup.iss
[Code]
#include "scripts\products.iss"

#include "scripts\products\stringversion.iss"
#include "scripts\products\winversion.iss"
#include "scripts\products\fileversion.iss"
#include "scripts\products\dotnetfxversion.iss"

#include "scripts\products\dotnetfx46.iss"
#include "scripts\products\msiproduct.iss" ; apprently vcredist2013 needs dist
#include "scripts\products\vcredist2013.iss"

function InitializeSetup(): boolean;
begin
	initwinversion();
    dotnetfx46(50); // min allowed version is 4.5.0
	SetForceX86(true); // force 32-bit install of next products
	vcredist2013();
	SetForceX86(false); // disable forced 32-bit install again
	Result := true;
end;
