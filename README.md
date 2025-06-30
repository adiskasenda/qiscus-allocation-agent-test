QISCUS ALLOCATION AGENT SERVICE

Requirement: 
- Go Chi
- Redis

How to Run?
1. create .env file
2. copy and paste from .example.env
3. fill the .env credential
4.  run with code bellow

```
go run main.go
```

Flow:
1. ERD Diagram for incoming chat to enqueue
```
flowchart TD
    A(["Start"]) --> B["Incoming request chat"]
    B --> C["Unmarshal request payload"]
    C --> n1["isEmpty roomId?"]
    n1 --> n2["Return error"] & n3["Enqueue to Redis"]
    B@{ shape: rect}
    n1@{ shape: diam}
```  

![image](https://github.com/user-attachments/assets/d04de6d6-e261-4c97-bf66-3ba4fac5dd61)

2. ERD Diagram for process allocation agent
```
flowchart TD
    A(["Start"]) --> B["Retry worker 30 sec"]
    B --> C["Get RoomId from Redis"]
    C --> n1["isEmty RoomId?"]
    n1 -- yes --> n2["sleep 60 sec &amp; retry"]
    n1 -- no --> n3["GetAssignedAgent from redis"]
    n2 --> A
    n3 --> n4["isAssign?"]
    n4 -- yes --> n5["skip agent &amp; room"]
    n4 -- no --> n6["GetAvailableAgent"]
    n6 --> n12["isSyncCount ?"]
    n7["count &gt; 2 ?"] -- yes --> n8["skip agent"]
    n7 -- no --> n9["count = 0 ?"]
    n9 -- no --> n10["assign agent to room"]
    n9 -- yes --> n11["assign 2 room to agent"]
    n12 -- yes --> n7
    n12 -- no --> n13["sycn count"]
    n13 --> n7
    n11 --> n14["increment count"]
    n10 --> n14
    n14 --> n15(["end"])
    B@{ shape: rect}
    C@{ shape: rect}
    n1@{ shape: diam}
    n4@{ shape: diam}
    n7@{ shape: diam}
    n9@{ shape: diam}
```

![image](https://github.com/user-attachments/assets/827a2d14-f711-4393-839d-e670df5c80d1)
  
3. Sequence Diagram
```
sequenceDiagram
  participant P1 as Customer
  participant P2 as Omnichannel
  participant P3 as Service Allocation Agent
  participant P4 as Redis
  participant P5 as Qiscus Api

  P3 ->> P4: Init connection
  P4 ->> P3: success connection
  P1 ->> P2: New Message (wa, fb, ig, etc)
  P2 ->> P3: Request Message via webhook
  P3 ->> P4: enqueue roomId
  P4 ->> P3: get roomId
  P4 ->> P3: get assignedAgent
  P3 ->> P5: get availableAgent
  P4 ->> P3: get agentChatCount
  P3 ->> P5: assignAgent
  P3 ->> P4: set assignAgent & increment chat count
  P3 ->> P4: enqueue roomId
  P5 ->> P2: chat served and agent filled
  P5 ->> P1: customer get agent in room chat
  P2 ->> P5: resolve chat
  P5 ->> P1: finish chat
```

![image](https://github.com/user-attachments/assets/6079bc59-0d4c-481f-a329-b2cfd016994f)
