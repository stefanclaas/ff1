# ff1
NIST FF1 CLI tool. Uses input files consisting of one large hex string.

Key creation: $ openssl rand -hex 32 > key.
Tweak creation: $ openssl rand -hex 8 > tweak.

Input message must consist of one (large) hex hex line, created with xxd,
or my base16 encoder/decoder, with line length set to 0.
