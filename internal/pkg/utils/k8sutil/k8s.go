package k8sutil

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	genericoptions "wechatbot/internal/pkg/options"
	"wechatbot/internal/pkg/utils/barkutil"
	"wechatbot/pkg/log"

	k8scorev1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func GetConfig() (opt *genericoptions.K8SOptions, err error) {
	// 获取namespace
	nsByte, errRead := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if errRead != nil {
		err = errRead
		return
	}
	namespace := string(nsByte)
	if namespace == "" {
		err = errors.New("no namespace")
		return
	}
	config, errCfg := rest.InClusterConfig()
	if errCfg != nil {
		err = errCfg
		return
	}
	// creates the clientset
	clientset, errClient := kubernetes.NewForConfig(config)
	if err != nil {
		err = errClient
		return
	}
	if clientset == nil || err != nil {
		if err == nil {
			err = errors.New("no k8s clientset")
		}
		return
	}
	HOSTNAME := os.Getenv("HOSTNAME")
	opt = &genericoptions.K8SOptions{
		Clientset: clientset,
		PodName:   HOSTNAME,
		Namespace: namespace,
	}
	return
}

func GetPods(clientset *kubernetes.Clientset, ns string) (*k8scorev1.PodList, error) {
	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), k8smetav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func GetPodInfo(clientset *kubernetes.Clientset, ns string, name string) (*k8scorev1.Pod, error) {
	// 获取Pod的状态
	pod, err := clientset.CoreV1().Pods(ns).Get(context.TODO(), name, k8smetav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func IsPodReady(pod *k8scorev1.Pod) (isReady bool) {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == k8scorev1.PodReady && condition.Status == k8scorev1.ConditionTrue {
			isReady = true
			break
		}
	}
	return
}

type LeaseEventRecorder struct {
	Name string
}

func (s *LeaseEventRecorder) Eventf(obj runtime.Object, eventType, reason, message string, args ...interface{}) {
	log.Infof("eventType:%s reason:%s message:%s args:%v", eventType, reason, message, args)
	if eventType != "Normal" {
		barkutil.SendMsg(fmt.Sprintf("【ms-go】 leaderelection EventRecorder"), fmt.Sprintf("eventType:%s reason:%s message:%s", eventType, reason, message))
	}
}

func StartLease(opt *genericoptions.K8SOptions, callbacks leaderelection.LeaderCallbacks, ctx context.Context) {
	// 指定锁的资源对象，这里使用了Lease资源，还支持configmap，endpoint，或者multilock(即多种配合使用)
	lock := &resourcelock.LeaseLock{
		LeaseMeta: k8smetav1.ObjectMeta{
			Name:      opt.LeaseName,
			Namespace: opt.Namespace,
		},
		Client: opt.Clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: opt.PodName,
			EventRecorder: &LeaseEventRecorder{
				Name: opt.PodName,
			},
		},
	}
	healthzAdaptor := leaderelection.NewLeaderHealthzAdaptor(10 * time.Second)
	// 定义Leader选举配置
	leaderConfig := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		LeaseDuration:   15 * time.Second, // Lease持续时间
		RenewDeadline:   10 * time.Second, // 更新截止时间
		RetryPeriod:     2 * time.Second,  // 重试间隔
		Callbacks:       callbacks,
		ReleaseOnCancel: true,
		WatchDog:        healthzAdaptor,
	}
	// 开始Leader选举
	leaderelection.RunOrDie(ctx, leaderConfig)
}
