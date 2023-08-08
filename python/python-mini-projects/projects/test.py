import requests

API_URL = "https://api-inference.huggingface.co/models/bigscience/bloom-7b1"
headers = {"Authorization": "Bearer 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"}

def query(payload):
	response = requests.post(API_URL, headers=headers, json=payload)
	return response.json()
	
output = query({
	"inputs": "Can you please let us know more details about your ",
})