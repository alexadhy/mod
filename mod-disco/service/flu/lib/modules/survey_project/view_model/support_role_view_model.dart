import 'package:flutter/material.dart';
import 'package:flutter_modular/flutter_modular.dart';
import 'package:mod_disco/core/shared_repositories/survey_project_repo.dart';
import 'package:mod_disco/core/shared_repositories/survey_user_repo.dart';
import 'package:mod_disco/core/shared_services/dynamic_widget_service.dart';
import 'package:mod_disco/rpc/v2/mod_disco_models.pb.dart';
import 'package:sys_core/sys_core.dart';
import 'package:sys_share_sys_account_service/pkg/shared_repositories/auth_repo.dart';
import 'package:sys_share_sys_account_service/pkg/shared_services/base_model.dart';
import 'package:sys_share_sys_account_service/sys_share_sys_account_service.dart';
import 'package:sys_share_sys_account_service/view/widgets/view_model/auth_nav_view_model.dart';

class SupportRoleViewModel extends BaseModel {
  Project _project;
  String _accountId = "";
  List<SurveyProject> _surveyProjects;
  NewSurveyUserRequest _nsuReq;
  List<List<SupportRoleType>> _srtLists = [];
  List<SupportRoleType> _srtList = [];
  Map<String, double> _minHours = {};
  DynamicWidgetService dwService = DynamicWidgetService();
  bool _isLoading = false;
  bool _isLoggedOn = false;

  Project get project => _project;

  NewSurveyUserRequest get nsuReq => _nsuReq;

  bool get isLoading => _isLoading;

  bool get isLoggedOn => _isLoggedOn;

  List<SupportRoleType> get supportRoles => _srtList;

  Map<String, double> get minHours => _minHours;
  Map<String, NewSupportRoleValue> _supportRoleMap = {};

  void setLoading(bool value) {
    _isLoading = value;
    notifyListeners();
  }

  // Constructor
  SupportRoleViewModel(
      {@required Project project,
      @required String accountId,
      @required NewSurveyUserRequest newSurveyUserRequest,
      @required List<SurveyProject> surveyProjectList}) {
    _project = project;
    _accountId = accountId;
    _nsuReq = newSurveyUserRequest;
    _surveyProjects = surveyProjectList;
  }

  void _isLoggedIn() {
    final isLoggedOn = Modular.get<AuthNavViewModel>().isLoggedIn;
    _isLoggedOn = isLoggedOn;
    notifyListeners();
  }

  void isUserLoggedIn() {
    return _isLoggedIn();
  }

  // init
  Future<void> initOnReady() async {
    setLoading(true);
    _surveyProjects.forEach((element) {
      _srtLists.add(element.supportRoleTypes);
    });
    _srtList = _srtLists.expand((i) => i).toList();
    isUserLoggedIn();
    notifyListeners();
    setLoading(false);
  }

  void selectMinHours(double value, String id) {
    _minHours[id] = value;
    _supportRoleMap[id] = SurveyProjectRepo.createSupportRoleValue(
      pledged: value.toInt(),
      surveyUserRefName: _nsuReq.surveyUserName,
      supportRoleTypeRefId: id,
    );
    notifyListeners();
  }

  Future<void> onSave(BuildContext context) async {
    List<NewSupportRoleValue> _srvList = [];
    _supportRoleMap.forEach((key, value) {
      _srvList.add(value);
    });
    final _userRole = UserRoles()
      ..role = Roles.USER
      ..projectId = _project.id
      ..orgId = _project.orgId;
    if (!isLoggedOn) {
      showDialog(
        context: context,
        builder: (context) => AuthDialog(
          isSignIn: false,
          userRole: _userRole,
          callback: () async {
            _accountId = getTempAccountId();
            _nsuReq.sysAccountUserRefId = _accountId;
            await SurveyUserRepo.newSurveyUser(
              surveyProjectId: _nsuReq.surveyProjectRefId,
              sysAccountAccountRefId: _nsuReq.sysAccountUserRefId,
              surveyUserName: _nsuReq.surveyUserName,
              userNeedsValueList: _nsuReq.userNeedValues,
              supportRoleValueList: _srvList,
            ).then((_) {
              notify(
                context: context,
                message:
                    "You've joined ${project.name}, login to see your detail",
                error: false,
              );
            });
          },
        ),
      );
    } else {
      final _su = await SurveyUserRepo.getSurveyUser(
        surveyProjectId: _nsuReq.surveyProjectRefId,
        sysAccountUserRefId: _nsuReq.sysAccountUserRefId,
      ).then((s) => s);
      if (_su == null) {
        await SurveyUserRepo.newSurveyUser(
          surveyProjectId: _nsuReq.surveyProjectRefId,
          sysAccountAccountRefId: _nsuReq.sysAccountUserRefId,
          surveyUserName: _nsuReq.surveyUserName,
          userNeedsValueList: _nsuReq.userNeedValues,
          supportRoleValueList: _srvList,
        ).then((_) {
          notify(
            context: context,
            message: "You've joined ${project.name}",
            error: false,
          );
        }).catchError((e) {
          notify(
            context: context,
            message: "Error submitting survey user ${e.toString()}",
            error: true,
          );
        });
      } else {
        final _unvs =
            _convertUserNeedsValues(_su.userNeedValues, _nsuReq.userNeedValues);
        final _srvs = _convertSupportRoleValues(
            _su.supportRoleValues, _nsuReq.supportRoleValues);
        await SurveyUserRepo.updateSurveyUser(
          surveyUserId: _su.surveyUserId,
          userNeedsValueList: _unvs,
          supportRoleValueList: _srvs,
        )
            .then((_) => notify(
                context: context,
                message: "Successfully updated your survey value",
                error: false))
            .catchError((e) {
          notify(
            context: context,
            message: "error updating your survey value: ${e.toString()}",
            error: true,
          );
        });
      }
    }
  }
}

List<UserNeedsValue> _convertUserNeedsValues(
    List<UserNeedsValue> oldValues, List<NewUserNeedsValue> newValues) {
  List<UserNeedsValue> updatedValues =
      List<UserNeedsValue>.empty(growable: true);
  oldValues.forEach((oldVal) {
    newValues.forEach((newVal) {
      if (oldVal.userNeedsTypeRefId == newVal.userNeedsTypeRefId) {
        updatedValues.add(UserNeedsValue(
          id: oldVal.id,
          surveyUserRefId: oldVal.surveyUserRefId,
          userNeedsTypeRefId: newVal.userNeedsTypeRefId,
          comments: newVal.comments,
        ));
      }
    });
  });
  return updatedValues;
}

List<SupportRoleValue> _convertSupportRoleValues(
    List<SupportRoleValue> oldValues, List<NewSupportRoleValue> newValues) {
  List<SupportRoleValue> updatedValues =
      List<SupportRoleValue>.empty(growable: true);

  oldValues.forEach((oldVal) {
    newValues.forEach((newVal) {
      if (oldVal.supportRoleTypeRefId == newVal.supportRoleTypeRefId) {
        updatedValues.add(SupportRoleValue(
          id: oldVal.id,
          surveyUserRefId: oldVal.surveyUserRefId,
          supportRoleTypeRefId: newVal.supportRoleTypeRefId,
          pledged: newVal.pledged,
          comment: newVal.comment,
        ));
      }
    });
  });

  return updatedValues;
}
