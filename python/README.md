# Mista Verify — Python

Uses the official [`mista`](https://github.com/mista-io/mista-python) package (Python 3.9+).

```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
export MISTA_API_TOKEN="your_mista_api_token_here"
```

| File | What it shows |
| --- | --- |
| [`quickstart.py`](quickstart.py) | Send a code, type it in, check it |
| [`channels.py`](channels.py) | Pick a channel — e.g. WhatsApp with SMS fallback — and look up status |
| [`async_quickstart.py`](async_quickstart.py) | Same flow with `AsyncMista` |
| [`server.py`](server.py) | Flask API for sign-up / login verification |

```bash
python quickstart.py +15555550100
python channels.py +15555550100 --channel whatsapp_sms
python async_quickstart.py +15555550100
python server.py
```

`+15555550100` is a placeholder — use a real number you can receive messages on.

## The SDK in 10 lines

```python
from mista import Mista

mista = Mista()  # reads MISTA_API_TOKEN

verification = mista.verify.start(to="+15555550100", channel="auto")
result = mista.verify.check(sid=verification["sid"], code="123456")

if result["verified"]:
    ...  # phone number confirmed
else:
    print(result.get("reason"))  # invalid_code | expired | max_attempts | not_open
```
