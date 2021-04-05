package repo

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.amplifyedge.org/mod-v2/mod-disco/service/go/pkg/dao"
	sysCoreSvc "go.amplifyedge.org/sys-v2/sys-core/service/go/pkg/coredb"

	discoRpc "go.amplifyedge.org/mod-v2/mod-disco/service/go/rpc/v2"
	sharedAuth "go.amplifyedge.org/sys-share-v2/sys-account/service/go/pkg/shared"
)

func (md *ModDiscoRepo) NewSurveyProject(ctx context.Context, in *discoRpc.NewSurveyProjectRequest) (*discoRpc.SurveyProject, error) {
	if in == nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot insert survey project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	// make sure that the sys-account project exists
	exists, sysAccountProjectId, err := md.sysAccountProjectExists(ctx, in.SysAccountProjectRefId, in.SysAccountProjectRefName)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot insert survey project, non existent sys-account-project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	if !exists {
		return nil, status.Errorf(codes.InvalidArgument, "cannot insert survey project: non-existent sys-account-project", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	in.SysAccountProjectRefId = sysAccountProjectId
	allowed := md.allowDiscoProject(ctx, "", in.SysAccountProjectRefId)
	if !allowed {
		return nil, status.Errorf(codes.PermissionDenied, "cannot insert survey project: permission denied", sharedAuth.Error{Reason: sharedAuth.ErrInsufficientRights})
	}
	md.log.Debugf("SysAccountProjectId: %s", sysAccountProjectId)
	sp, err := md.store.InsertSurveyProject(in)
	if err != nil {
		return nil, err
	}
	return sp, nil
}

func (md *ModDiscoRepo) GetSurveyProject(ctx context.Context, in *discoRpc.IdRequest) (*discoRpc.SurveyProject, error) {
	if in == nil || (in.SurveyProjectId == "" && in.SysAccountProjectId == "") {
		return nil, status.Errorf(codes.InvalidArgument, "cannot get survey project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	params := map[string]interface{}{}
	if in.GetSurveyProjectId() != "" {
		params["survey_project_id"] = in.GetSurveyProjectId()
	}
	if in.GetSysAccountProjectId() != "" {
		params["sys_account_project_ref_id"] = in.GetSysAccountProjectId()
	}
	sp, err := md.store.GetSurveyProject(params)
	if err != nil {
		return nil, err
	}
	return md.store.ToRpcSurveyProject(sp)
}

func (md *ModDiscoRepo) ListSurveyProject(ctx context.Context, in *discoRpc.ListRequest) (*discoRpc.ListResponse, error) {
	if in == nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot list survey project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	var cursor int64
	orderBy := in.OrderBy
	var err error
	filter := map[string]interface{}{}
	if in.Filters != nil && len(in.Filters) > 0 {
		filter, err = sysCoreSvc.UnmarshalToMap(in.Filters)
		if err != nil {
			return nil, err
		}
	}
	if in.IdRequest != nil {
		if in.IdRequest.SurveyProjectId != "" {
			filter["survey_project_id"] = in.IdRequest.SurveyProjectId
		}
		if in.IdRequest.SysAccountProjectId != "" {
			filter["sys_account_project_ref_id"] = in.IdRequest.SysAccountProjectId
		}
	}
	cursor, err = md.getCursor(in.GetCurrentPageId())
	if err != nil {
		return nil, err
	}
	if in.GetIsDescending() {
		orderBy += " DESC"
	} else {
		orderBy += " ASC"
	}
	limit := in.PerPageEntries
	if limit == 0 {
		limit = dao.DefaultLimit
	}
	daoSurveyProjects, next, err := md.store.ListSurveyProject(filter, orderBy, limit, cursor, in.GetMatcher())
	if err != nil {
		return nil, err
	}
	var pkgSurveyProjects []*discoRpc.SurveyProject
	for _, sp := range daoSurveyProjects {
		surveyProject, err := md.store.ToRpcSurveyProject(sp)
		if err != nil {
			return nil, err
		}
		pkgSurveyProjects = append(pkgSurveyProjects, surveyProject)
	}
	return &discoRpc.ListResponse{
		SurveyProjects: pkgSurveyProjects,
		NextPageId:     *next,
	}, nil
}

func (md *ModDiscoRepo) UpdateSurveyProject(ctx context.Context, in *discoRpc.UpdateSurveyProjectRequest) (*discoRpc.SurveyProject, error) {
	if in == nil || in.SurveyProjectId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cannot update survey project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	daoSurveyProj, err := md.store.GetSurveyProject(map[string]interface{}{"survey_project_id": in.SurveyProjectId})
	if err != nil {
		return nil, err
	}
	allowed := md.allowDiscoProject(ctx, "", daoSurveyProj.SysAccountProjectRefId)
	if !allowed {
		return nil, status.Errorf(codes.PermissionDenied, "cannot insert disco project: permission denied", sharedAuth.Error{Reason: sharedAuth.ErrInsufficientRights})
	}
	if err := md.store.UpdateSurveyProject(in); err != nil {
		return nil, err
	}
	daoSurveyProj, err = md.store.GetSurveyProject(map[string]interface{}{"survey_project_id": in.SurveyProjectId})
	if err != nil {
		return nil, err
	}
	return md.store.ToRpcSurveyProject(daoSurveyProj)
}

func (md *ModDiscoRepo) DeleteSurveyProject(ctx context.Context, in *discoRpc.IdRequest) (*emptypb.Empty, error) {
	if in == nil || in.SurveyProjectId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cannot update survey project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	daoSurveyProj, err := md.store.GetSurveyProject(map[string]interface{}{"survey_project_id": in.SurveyProjectId})
	if err != nil {
		return nil, err
	}
	allowed := md.allowDiscoProject(ctx, "", daoSurveyProj.SysAccountProjectRefId)
	if !allowed {
		return nil, status.Errorf(codes.PermissionDenied, "cannot insert disco project: permission denied", sharedAuth.Error{Reason: sharedAuth.ErrInsufficientRights})
	}
	err = md.store.DeleteSurveyProject(in.SurveyProjectId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
