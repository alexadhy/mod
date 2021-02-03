package repo

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/amplify-cms/mod/mod-disco/service/go/pkg/dao"
	discoRpc "github.com/amplify-cms/mod/mod-disco/service/go/rpc/v2"
	sharedAuth "github.com/amplify-cms/sys-share/sys-account/service/go/pkg/shared"
	sharedConfig "github.com/amplify-cms/sys-share/sys-core/service/config"
	sysCoreSvc "github.com/amplify-cms/sys/sys-core/service/go/pkg/coredb"
)

func (md *ModDiscoRepo) NewDiscoProject(ctx context.Context, in *discoRpc.NewDiscoProjectRequest) (*discoRpc.DiscoProject, error) {
	if in == nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot insert disco project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	// make sure that the sys-account project exists
	exists, sysAccountProjectId, err := md.sysAccountProjectExists(ctx, in.SysAccountProjectRefId, in.SysAccountProjectRefName)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot insert disco project, non existent sys-account-project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	if !exists {
		return nil, status.Errorf(codes.InvalidArgument, "cannot insert disco project: non-existent sys-account-project", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	in.SysAccountProjectRefId = sysAccountProjectId
	var imgResourceIds []string
	if len(in.GetImageFilepath()) != 0 {
		for _, imgPath := range in.GetImageFilepath() {
			resId, err := md.frepo.UploadFile(imgPath, nil)
			if err != nil {
				return nil, err
			}
			imgResourceIds = append(imgResourceIds, resId.ResourceId)
		}
	}
	if len(in.GetImageUploadArrays()) != 0 {
		for _, imgData := range in.GetImageUploadArrays() {
			imgBytes, err := sharedConfig.DecodeB64(imgData)
			if err != nil {
				return nil, err
			}
			resId, err := md.frepo.UploadFile("", imgBytes)
			if err != nil {
				return nil, err
			}
			imgResourceIds = append(imgResourceIds, resId.ResourceId)
		}
	}
	dp, err := md.store.InsertDiscoProject(in, imgResourceIds)
	if err != nil {
		return nil, err
	}
	return md.fetchProjectImages(dp)
}

func (md *ModDiscoRepo) fetchProjectImages(dp *discoRpc.DiscoProject) (*discoRpc.DiscoProject, error) {
	var imgData [][]byte
	for _, resId := range dp.ImageResourceIds {
		img, err := md.frepo.DownloadFile("", resId)
		if err != nil {
			return nil, err
		}
		imgData = append(imgData, img.Binary)
	}
	dp.ProjectImages = imgData
	return dp, nil
}

func (md *ModDiscoRepo) GetDiscoProject(ctx context.Context, in *discoRpc.IdRequest) (*discoRpc.DiscoProject, error) {
	if in == nil || (in.GetDiscoProjectId() == "" && in.GetSysAccountProjectId() == "") {
		return nil, status.Errorf(codes.InvalidArgument, "cannot get disco project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	params := map[string]interface{}{}
	if in.GetDiscoProjectId() != "" {
		params["project_id"] = in.GetDiscoProjectId()
	}
	if in.GetSysAccountProjectId() != "" {
		params["sys_account_project_ref_id"] = in.GetSysAccountProjectId()
	}
	dp, err := md.store.GetDiscoProject(params)
	if err != nil {
		return nil, err
	}
	dproj, err := dp.ToPkgDiscoProject()
	if err != nil {
		return nil, err
	}
	return md.fetchProjectImages(dproj)
}

func (md *ModDiscoRepo) ListDiscoProject(ctx context.Context, in *discoRpc.ListRequest) (*discoRpc.ListResponse, error) {
	if in == nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot list survey user: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
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
	daoDiscoProjects, next, err := md.store.ListDiscoProject(filter, orderBy, limit, cursor, in.GetMatcher())
	if err != nil {
		return nil, err
	}
	var pkgDiscoProjects []*discoRpc.DiscoProject
	for _, su := range daoDiscoProjects {
		dp, err := su.ToPkgDiscoProject()
		if err != nil {
			return nil, err
		}
		dp, err = md.fetchProjectImages(dp)
		if err != nil {
			return nil, err
		}
		pkgDiscoProjects = append(pkgDiscoProjects, dp)
	}
	return &discoRpc.ListResponse{
		DiscoProjects: pkgDiscoProjects,
		NextPageId:    *next,
	}, nil
}

func (md *ModDiscoRepo) UpdateDiscoProject(ctx context.Context, in *discoRpc.UpdateDiscoProjectRequest) (*discoRpc.DiscoProject, error) {
	if in == nil || in.ProjectId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cannot update disco project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	if in.GetImageUploads() != nil && len(in.GetImageUploads()) != 0 {
		for _, imgBytes := range in.GetImageUploads() {
			res, err := md.frepo.UploadFile("", imgBytes)
			if err != nil {
				return nil, err
			}
			in.ImageResourceIds = append(in.ImageResourceIds, res.ResourceId)
		}
	}
	if err := md.store.UpdateDiscoProject(in); err != nil {
		return nil, err
	}
	daoDiscoProject, err := md.store.GetDiscoProject(map[string]interface{}{"project_id": in.ProjectId})
	if err != nil {
		return nil, err
	}
	dp, err := daoDiscoProject.ToPkgDiscoProject()
	if err != nil {
		return nil, err
	}
	return md.fetchProjectImages(dp)
}

func (md *ModDiscoRepo) DeleteDiscoProject(ctx context.Context, in *discoRpc.IdRequest) (*emptypb.Empty, error) {
	if in == nil || (in.DiscoProjectId == "" && in.SysAccountProjectId == "" && in.SysAccountOrgId == "") {
		return nil, status.Errorf(codes.InvalidArgument, "cannot delete disco project: %v", sharedAuth.Error{Reason: sharedAuth.ErrInvalidParameters})
	}
	err := md.store.DeleteDiscoProject(in.DiscoProjectId, in.SysAccountProjectId, in.SysAccountOrgId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
