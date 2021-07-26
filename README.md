## MQTT for Hikvision camera

Capture a frame every 5 seconds and publish to mqtt.

Run:
```LD_LIBRARY_PATH=/hiksdk/lib/ ./hikmqtt -c hikmqtt.json```

Config hikmqtt.json:
```
{
    "Url": "mqtt://127.0.0.1:1883",
    "Username": "user",
    "Password": "password",
    "Cams": [
        {
            "Ip": "192.168.0.1",
            "Username": "camuser",
            "Password": "campassword",
            "Name": "First cam"
        },
        {
            "Ip": "192.168.0.2",
            "Username": "camuser1",
            "Password": "campassword1",
            "Name": "Second cam"
        }
    ]
}
```
