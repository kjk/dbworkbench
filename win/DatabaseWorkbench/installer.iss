; Requires https://www.dropbox.com/s/lkr9qh3uj0hkqqp/idpsetup-1.5.0.exe?dl=0 to be installed
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

#include <idp.iss>

#define MyAppName "Database Workbench"
#define MyAppVersion "9.9"
#define MyAppPublisher "Krzysztof Kowalczyk"
#define MyAppURL "http://www.databaseworkbench.com"
#define MyAppExeName "DatabaseWorkbench.exe"

[Setup]
; NOTE: The value of AppId uniquely identifies this application.
; Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={784A5D3A-FB9F-4E09-809C-40F639F408D6}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
;AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={localappdata}\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
OutputDir=bin\Release
OutputBaseFilename=DatabaseWorkbench-setup-inno
SetupIconFile=icon.ico
Compression=lzma
SolidCompression=yes
; Vista SP 1, http://www.jrsoftware.org/ishelp/index.php?topic=winvernotes
; http://www.gaijin.at/en/lstwinver.php
MinVersion=6.0.6002
DisableDirPage=yes
;DisableWelcomePage=yes
;DisableReadyPage=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"

[Files]
Source: "bin\Release\DatabaseWorkbench.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\DatabaseWorkbench.exe.config"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\dbworkbench.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\dbworkbench.dat"; DestDir: "{app}"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{commondesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent

[Code]
function Framework45IsNotInstalled(): Boolean;
var
  bSuccess: Boolean;
  regVersion: Cardinal;
begin
  Result := True;

  bSuccess := RegQueryDWordValue(HKLM, 'Software\Microsoft\NET Framework Setup\NDP\v4\Full', 'Release', regVersion);
  if (True = bSuccess) and (regVersion >= 378389) then begin
    Result := False;
  end;
end;

procedure InitializeWizard;
begin
  if Framework45IsNotInstalled() then
  begin
    idpAddFile('http://go.microsoft.com/fwlink/?LinkId=397707', ExpandConstant('{tmp}\NetFrameworkInstaller.exe'));
    idpDownloadAfter(wpReady);
  end;
end;

procedure InstallFramework;
var
  StatusText: string;
  ResultCode: Integer;
begin
  StatusText := WizardForm.StatusLabel.Caption;
  WizardForm.StatusLabel.Caption := 'Installing .NET Framework 4.5.2. This might take a few minutes...';
  WizardForm.ProgressGauge.Style := npbstMarquee;
  try
    if not Exec(ExpandConstant('{tmp}\NetFrameworkInstaller.exe'), '/passive /norestart', '', SW_SHOW, ewWaitUntilTerminated, ResultCode) then
    begin
      MsgBox('.NET installation failed with code: ' + IntToStr(ResultCode) + '.', mbError, MB_OK);
    end;
  finally
    WizardForm.StatusLabel.Caption := StatusText;
    WizardForm.ProgressGauge.Style := npbstNormal;

    DeleteFile(ExpandConstant('{tmp}\NetFrameworkInstaller.exe'));
  end;
end;


procedure CurStepChanged(CurStep: TSetupStep);
begin
  case CurStep of
    ssPostInstall:
      begin
        if Framework45IsNotInstalled() then
        begin
          InstallFramework();
        end;
      end;
  end;
end;
