<%@page import="javax.servlet.ServletInputStream,javax.servlet.http.HttpServletRequest"%>
<%@page import="java.io.BufferedReader,java.io.InputStreamReader"%>
<%@page import="java.io.PrintWriter"%>
<%@page import="java.util.Enumeration"%>
<%@page import="java.util.Map"%>
<%

      BufferedReader br = request.getReader();
      String prefix="payment_intimation";

      PrintWriter writer = new PrintWriter("yourfilenamewithpath", "UTF-8");
        String line = "";
        out.println("----META DATA-----");
        out.println("Remote Address:"+request.getRemoteAddr());
        out.println("Content Length:"+request.getContentLength());
        out.println("Content Type:"+request.getContentType());
        out.println("Character encoding:"+request.getCharacterEncoding());
        out.println("Auth Type:"+request.getAuthType());
        out.println("Context Path:"+request.getContextPath());
        out.println("Local Address:"+request.getLocalAddr());
        out.println("Local Name:"+request.getLocalName());
        out.println("Local Port:"+request.getLocalPort());
        out.println("Method:"+request.getMethod());
        out.println("Path Info:"+request.getPathInfo());
        out.println("Path Transalated:"+request.getPathTranslated());
        out.println("Protocol:"+request.getProtocol());
        out.println("QueryString:"+request.getQueryString());
        out.println("Remote Host:"+request.getRemoteHost());
        out.println("Remote User:"+request.getRemoteUser());
        out.println("Requested Session:"+ request.getRequestedSessionId());
        out.println("Request URI:"+ request.getRequestURI());
        out.println("Request URL:"+ request.getRequestURL());
        out.println("Scheme:"+ request.getScheme());
        out.println("ServerName:"+ request.getServerName());
        out.println("ServerPort:"+ request.getServerPort());
        out.println("Servlet Path:"+ request.getServletPath());
        out.println("----HEADER-----");
        Enumeration<String> headerNames = request.getHeaderNames();
        while (headerNames.hasMoreElements()) {
             String headerName = headerNames.nextElement();
             Enumeration<String> headers = request.getHeaders(headerName);
           while (headers.hasMoreElements()) {
                  String headerValue = headers.nextElement();
                  out.println(headerName+":"+headerValue);
             }
        }
        out.println("----PARAMETERS-----");
        Map<String, String[]> parameters = request.getParameterMap();
        for(String parameter : parameters.keySet()) {
                String[] values = parameters.get(parameter);
                for (int i=0; i < values.length;i++) {
                  out.println(parameter+":"+values[i]);
                }
        }
out.println("----BODY-----");
        while((line = br.readLine()) != null) {
                out.println(line);
        }
     writer.close();
%>