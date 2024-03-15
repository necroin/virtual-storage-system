C:/Storage/"Program Files"/Git/usr\bin\openssl.exe req -x509 -newkey rsa:4096 -sha256 ^
-days 3650 -nodes -keyout vss.key -out vss.crt ^
-subj "/CN=vss" ^
-addext "subjectAltName=DNS:vss,DNS:localhost"