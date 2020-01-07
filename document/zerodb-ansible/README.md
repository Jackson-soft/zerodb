
# Install ansible   sshpass

brew install ansible

brew install sshpass
## 或者
brew install http://git.io/sshpass.rb

# build zerokeeper zeroproxy
make clean
make

# 定义参数

1. 修改你的用户名,密码和私钥位置
```
    ansible_ssh_user: 'youre user name'
    ansible_ssh_private_key_file: 'your private key path'
    ansible_ssh_pass: 'your password'
 ```

# daily 环境部署

```
ansible-playbook -i hosts -e "env=daily GOPATH=`echo $GOPATH`" zeroproxy.yml
ansible-playbook -i hosts -e "env=daily GOPATH=`echo $GOPATH`" zerokeeper.yml
```

# public 环境部署

```
ansible-playbook -i zerodb-ansible/hosts -e "env=public GOPATH=`echo $GOPATH`" zerodb-ansible/zerokeeper.yml
ansible-playbook -i zerodb-ansible/hosts -e "env=public GOPATH=`echo $GOPATH`" zerodb-ansible/zeroproxy.yml
```
