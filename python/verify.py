from mista import Mista

mista = Mista()  # reads MISTA_API_TOKEN

# 1. Send a code (WhatsApp first, SMS fallback)
verification = mista.verify.start(to="+15555550100", channel="whatsapp_sms")

# 2. Ask the user for it
code = input("Code: ")

# 3. Check it
result = mista.verify.check(sid=verification["sid"], code=code)
print("Verified" if result["verified"] else f"Not verified: {result['reason']}")
