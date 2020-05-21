# doubtnut
Email Service


## Installation

```
OS Environment - Linux (Ubuntu)
Make Sure Go and dep is Installed in the system
Install dependencies 
cmd - dep ensure -v
./build.sh will built Go binary 
Run the binary ./doubtnut to start the server
Uses SendGrid to send the emails,Currently email are hardcoded but can be made dynamic based on userid
Add SendGrid Api key in .env
Generated pdf filename is questions.pdf

```
Payload 
```
curl -X POST \
  http://127.0.0.1:6000/api/v1/send/1212 \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "1": "As light from a star spreads out and weakens, do gaps form between the photons?",
    "2": "Can momentum be hidden to human eyes like how kinetic energy can be hidden as heat?",
    "3":"Can you make a shock wave of light by breaking the light barrier just like supersonic airplanes break the sound barrier?",
    "4":"How bad would a person'\''s injuries be if an elevator'\''s cables snapped at the 100th floor so that the elevator fell to the bottom?",
    "5":"How does a microwave oven heat up food even though it emits no thermal radiation?",
    "6":"If I'\''m on an elevator that breaks loose and plummets down the shaft, can I avoid harm by jumping at the last second?",
    "7":"Why were electrons chosen to be negatively charged? Wouldn'\''t it make more sense to call electrons positively charged because when they move they make electricity?"
  
    
    
}'




```
