# Quotify

Internal quote bot for quotes twice a day.

Note: will only use quotes of the format `@person "I actually said this"` any references to an `@user` in the quote itself will resolve to the full name.

### Deployment 
Host as two AWS lambda functions generation & delivery (create these via ''dep ensure & 'go generate')

Set environment vars for each lamdba function

##### Generation
- SLACK_TOKEN (Your slack token)
- SLACK_QUOTE_CHANNEL_ID (The ID (Not name) of the channel to read quotes from)

##### Delivery
- SLACK_CHANNEL (The name of the person or channel to recieve the rendered quotes)
- SLACK_WEBHOOK_URL (your slack webhook URL)
- IMAGE_GEN_ENDPOINT (endpoint of your url image renderer: Note this project uses the Lucid Cube URL Grab service)
- ENT_DATA (JSON encoded string array used)
