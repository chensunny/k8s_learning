## 简单部署一个 k8s 实例


### 相关资料




#### 编辑 镜像


```
docker build .  --tag  chensunny/k8slearning:0.1

docker run -it --rm -p 8080:8080  chensunny/k8slearning:0.1
```


#### 推到镜像仓库

```

docker login
docker push   chensunny/k8slearning:0.1

```

https://www.qikqiak.com/post/write-kubernets-golang-service-step-by-step/




#### 登录部署到k8s



gcloud container clusters get-credentials standard-cluster-1 --zone us-central1-a --project aftership-team-brandy



kubectl create ns crs

```
kubectl create deployment nginx --image=nginx
kubectl create service nodeport nginx --tcp 80:80
```



```
kubectl create deployment k8slearning --image=chensunny/k8slearning:0.1
kubectl create service loadbalancer k8slearning --tcp 80:8080

```


#### VSCode 创建 k8s 资源


https://mp.weixin.qq.com/s/eps5pkocUwlWcZhB3fdJIQ




