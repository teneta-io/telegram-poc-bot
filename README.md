# TENETA Proof of Concept Telegram bot

## Starting to communicate with a bot

```
/start
```

Responce:


Welcome to TENETA Concept Web3 World.

We have created a wallet for your:

Your wallet address: 'Provider_wallet'

**You can deploy your own Virtual Machine on TENETA for affordable price**

or

**You can run someone's Virtual Machine and earn some goods**


Read more:

[TENETA site](https://teneta.io)

[TENETA on GitHub](https://github.com/teneta-io)

[Contribute TENETA](mailto:join@teneta.io)


## For users who want to be a provider, and earn some goods providing computing power

```
/act as provider, and deploy someone’s VM
  /vCPU limit: 1, 2, 4, 8, 16, 32, 64, 128
  /RAM limit, GB: 0.5, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024
  /Storage limit, GB: 20, 40, 80, 120, 160, 200, 240, 300, 1000, 2000
  /network throughput limit, Mb: 100, 1000, 10000
  /open ports: [UDP:53, TCP:80, TCP:443]
```

Responce:

```
  “You are registered as a provider with limits:
  16 vCPU,
  RAM 32GB,
  Storage 120GB,
  Network 1000Mb,
  wait for tasks…”

```
Task proposals:
```
/Task proposal 1:
    ‘Task_id’
    1vCPU,
    RAM 2GB,
    Storage 20GB,
    Network 100Mb,
    Ubuntu 22.04 LTS,
    ‘ssh-rsa SSH_PUB_KEY’,
    Price = 1 coin/day

/Task proposal 2:
    ‘Task_id’
    2vCPU, RAM 4GB,
    Storage 40GB,
    Network 100Mb,
    Ubuntu 22.04 LTS,
    ‘ssh-rsa SSH_PUB_KEY’,
    Price = 2 coin/day

/Task proposal 3:
    ‘Task_id’
    4vCPU,
    RAM 8GB,
    Storage 80GB,
    Network 1000Mb,
    Ubuntu 22.04 LTS,
    ‘ssh-rsa SSH_PUB_KEY’,
    Price 8 coin/day
```
Responce:
```
You have commit to run
    ‘Task_id’
    2vCPU, RAM 4GB,
    Storage 40GB,
    Network 100Mb,
    Ubuntu 22.04 LTS,
    ‘ssh-rsa SSH_PUB_KEY’,
    Price = 2 coin/day”


/provide ssh access string:
  ...
```

### Get my commitment:

```
/get my commitment
```

Responce:

```
‘Task_id’ : {
  2vCPU,
  RAM 4GB,
  Storage 40GB,
  Network 100Mb,
  Ubuntu 22.04 LTS,
  ‘ssh-rsa SSH_PUB_KEY’, Price 1 coin/day},
Status : Executing,
Access : root@public.ip -p PORT_NUMBER,
Provider wallet: Provider_wallet
```

## For users who want to run some Virtual Machine somewhere for an affordable price

```
/act as customer, run your own VM
  /create task
  /how many vCPU: 1, 2, 4, 8, 16, 32, 64
  /how many RAM, GB: 0.5, 1, 2, 4, 8, 16, 32, 64, 128
	/how many storage, GB: 20, 40, 80, 120, 160, 200, 240, 300
	/network throughput, Mb: 100, 1000		
  /operation system: Ubuntu 22.04 LTS, Centos 7
	/open ports: TCP 80, UDP 443, UDP 1194
  /post your public key
	/price for day, coins
```

Responce:

```
“The task ‘task_id’ was published, wait for execution and getting access”

“Your task ‘task_id’ has been started.

Access: ```root@public.ip -p PORT_NUMBER```.

Provider wallet: Provider_wallet”
```

### Check existing tasks:

```
/get my tasks
```
Responce:
```
“‘Task_id’ : {1vCPU, RAM 2GB, Storage 20GB, Network 100Mb, Ubuntu 22.04 LTS,
‘ssh-rsa SSH_PUB_KEY’, Price 1 coin/day},
Status : Initiated

“‘Task_id’ : {2vCPU, RAM 4GB, Storage 40GB, Network 100Mb, Ubuntu 22.04 LTS,
‘ssh-rsa SSH_PUB_KEY’, Price 1 coin/day},
Status : Executing.

'root@public.ip -p PORT_NUMBER'.

Provider wallet: Provider_wallet
```

### Delete task:

```
/delete task: ‘Task_id’
```

Responce:

```
“Task ‘Task_id’ has been deleted”
```

### Get balance

```
/get my balance
```
Responce:
```
“Your balance is equal 0 coins”
```

### Send coins

```
/send coins
```

### Buy coins

```
/buy coins
```

### Sale coins

```
/sale coins
```

### Coin faucet

```
/faucet coins
```
