package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/zxnlx/common"
	"github.com/zxnlx/svc/proto/svc"
	"github.com/zxnlx/svc_api/plugin/form"
	"github.com/zxnlx/svc_api/proto/svc_api"
	"strconv"
)

type SvcApi struct {
	SvcService svc.SvcService
}

func (s *SvcApi) FindSvcById(ctx context.Context, req *svc_api.Request, resp *svc_api.Response) error {
	log.Info("接受 svcApi.FindSvcById 的请求")
	if _, ok := req.Get["svc_id"]; !ok {
		resp.StatusCode = 500
		return errors.New("参数异常")
	}
	//获取 svcId 参数
	svcIdString := req.Get["svc_id"].Values[0]
	svcId, err := strconv.ParseInt(svcIdString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}
	//获取svc相关信息
	svcInfo, err := s.SvcService.FindSvcByID(ctx, &svc.SvcId{
		Id: svcId,
	})
	if err != nil {
		common.Error(err)
		return err
	}
	//json 返回svc信息

	resp.StatusCode = 200
	b, _ := json.Marshal(svcInfo)
	resp.Body = string(b)
	return nil
}

func (s *SvcApi) AddSvc(ctx context.Context, req *svc_api.Request, resp *svc_api.Response) error {
	log.Info("添加svc服务")
	//处理port
	addSvcInfo := &svc.SvcInfo{}
	svcType, ok := req.Post["svc_type"]
	if ok && len(svcType.Values) > 0 {
		svcPort := &svc.SvcPort{}
		switch svcType.Values[0] {
		case "ClusterIP":
			port, err := strconv.ParseInt(req.Post["svc_port"].Values[0], 10, 32)
			if err != nil {
				common.Error(err)
				return err
			}
			svcPort.SvcPort = int32(port)
			targetPort, err := strconv.ParseInt(req.Post["svc_target_port"].Values[0], 10, 32)
			if err != nil {
				common.Error(err)
				return err
			}
			svcPort.SvcTargetPort = int32(targetPort)
			svcPort.SvcPortProtocol = req.Post["svc_port_protocol"].Values[0]
			addSvcInfo.SvcPort = append(addSvcInfo.SvcPort, svcPort)
		default:
			return errors.New("类型不支持")
		}
	}
	//form 类型转换到结构体中
	form.FromToSvcStruct(req.Post, addSvcInfo)
	response, err := s.SvcService.AddSvc(ctx, addSvcInfo)
	if err != nil {
		common.Error(err)
		return err
	}
	resp.StatusCode = 200
	b, _ := json.Marshal(response)
	resp.Body = string(b)
	return nil
}

func (s *SvcApi) DeleteSvcById(ctx context.Context, req *svc_api.Request, resp *svc_api.Response) error {
	log.Info("删除svc服务")
	if _, ok := req.Get["svc_id"]; !ok {
		return errors.New("参数异常")
	}
	//获取需要删除的ID
	svcIdString := req.Get["svc_id"].Values[0]
	svcId, err := strconv.ParseInt(svcIdString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}
	//调用后端服务删除
	response, err := s.SvcService.DeleteSvc(ctx, &svc.SvcId{
		Id: svcId,
	})
	if err != nil {
		common.Error(err)
		return err
	}
	resp.StatusCode = 200
	b, _ := json.Marshal(response)
	resp.Body = string(b)
	return nil
}

func (s *SvcApi) UpdateSvc(ctx context.Context, req *svc_api.Request, resp *svc_api.Response) error {
	log.Info("Received svcApi.UpdateSvc request")
	resp.StatusCode = 200
	b, _ := json.Marshal("{success:'成功访问/svcApi/UpdateSvc'}")
	resp.Body = string(b)
	return nil
}

func (s *SvcApi) Call(ctx context.Context, req *svc_api.Request, resp *svc_api.Response) error {
	log.Info("查询所有svc服务")
	allSvc, err := s.SvcService.FindAllSvc(ctx, &svc.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}
	resp.StatusCode = 200
	b, _ := json.Marshal(allSvc)
	resp.Body = string(b)
	return nil
}
