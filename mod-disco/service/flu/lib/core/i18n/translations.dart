import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class Translations {
  Locale locale;

  // ORGS

  String category() {
    return Intl.message(
      'Category',
      name: 'category',
      desc: 'Category',
      locale: locale.toString(),
    );
  }

  String campaignName() {
    return Intl.message(
      'Campaign Name',
      name: 'campaignName',
      desc: 'campaignName',
      locale: locale.toString(),
    );
  }

  String noCampaigns() {
    return Intl.message(
      'No Campaigns',
      name: 'noCampaigns',
      desc: 'No Campaigns',
      locale: locale.toString(),
    );
  }

  String notFound() {
    return Intl.message(
      'Not found',
      name: 'notFound',
      desc: 'Not Found',
      locale: locale.toString(),
    );
  }

  String actionType() {
    return Intl.message(
      'Type of Action',
      name: 'actionType',
      desc: 'Type of Action',
      locale: locale.toString(),
    );
  }

  String actionLocation() {
    return Intl.message(
      'Action Location',
      name: 'actionLocation',
      desc: 'Action Location',
      locale: locale.toString(),
    );
  }

  String time() {
    return Intl.message(
      'Time',
      name: 'time',
      desc: 'Time',
      locale: locale.toString(),
    );
  }

  String lengthOfTheAction() {
    return Intl.message(
      'Length of the Action',
      name: 'lengthOfTheAction',
      desc: 'Length of the Action',
      locale: locale.toString(),
    );
  }

  String goal() {
    return Intl.message(
      'Goal',
      name: 'goal',
      desc: 'Goal',
      locale: locale.toString(),
    );
  }

  String strategy() {
    return Intl.message(
      'Strategy',
      name: 'strategy',
      desc: 'strategy',
      locale: locale.toString(),
    );
  }

  String historicalPrecedents() {
    return Intl.message(
      'Historical Precedents',
      name: 'historicalPrecedents',
      desc: 'Historical Precedents',
      locale: locale.toString(),
    );
  }

  String backingEndorsingOrganizations() {
    return Intl.message(
      'Backing/Endorsing Organizations',
      name: 'backingEndorsingOrganizations',
      desc: 'Backing/Endorsing Organizations',
      locale: locale.toString(),
    );
  }

  String peopleAlreadyPledged() {
    return Intl.message(
      'People already pledged',
      name: 'peopleAlreadyPledged',
      desc: 'People already pledged',
      locale: locale.toString(),
    );
  }

  String weNeed() {
    return Intl.message(
      'We Need',
      name: 'weNeed',
      desc: 'We Need',
      locale: locale.toString(),
    );
  }

  String extrapolatedSimilarPastActions() {
    return Intl.message(
      'The following figures are extrapolated from similar past actions that both succeeded and failed',
      name: 'extrapolatedSimilarPastActions',
      desc:
          'The following figures are extrapolated from similar past actions that both succeeded and failed',
      locale: locale.toString(),
    );
  }

  String pioneersToStart() {
    return Intl.message(
      'Pioneers needed to start',
      name: 'pioneersToStart',
      desc: 'Pioneers needed to start',
      locale: locale.toString(),
    );
  }

  String rebelsMedia() {
    return Intl.message(
      'Rebels needed to trigger media',
      name: 'rebelsMedia',
      desc: 'Rebels needed to trigger media',
      locale: locale.toString(),
    );
  }

  String rebelsWin() {
    return Intl.message(
      'Rebels needed to win',
      name: 'rebelsWin',
      desc: 'Rebels needed to win',
      locale: locale.toString(),
    );
  }

  String contactDetails() {
    return Intl.message(
      'Contact Details',
      name: 'contactDetails',
      desc: 'Contact Details',
      locale: locale.toString(),
    );
  }

  String notReady() {
    return Intl.message(
      'Not Ready',
      name: 'notReady',
      desc: 'Not Ready',
      locale: locale.toString(),
    );
  }

  String ready() {
    return Intl.message(
      'Ready',
      name: 'ready',
      desc: 'Ready',
      locale: locale.toString(),
    );
  }

  String selectCampaign() {
    return Intl.message(
      'Select Campaign',
      name: 'selectCampaign',
      desc: 'Select Campaign',
      locale: locale.toString(),
    );
  }

  String campaignDetails() {
    return Intl.message(
      'Campaign Details',
      name: 'campaignDetails',
      desc: 'campaignDetails',
      locale: locale.toString(),
    );
  }

  String searchCampaigns() {
    return Intl.message(
      'Search Campaigns',
      name: 'searchCampaigns',
      desc: 'Search Campaigns',
      locale: locale.toString(),
    );
  }

  // USERNEEDS

  String yourNeeds() {
    return Intl.message(
      'Your Needs',
      name: 'yourNeeds',
      desc: 'Your Needs',
      locale: locale.toString(),
    );
  }

  String needsSatisifiedRequirement() {
    return Intl.message(
      'Please choose as many supports or needs you need satisfied to join the action.',
      name: 'needsSatisifiedRequirement',
      desc:
          'Please choose as many supports or needs you need satisfied to join the action.',
      locale: locale.toString(),
    );
  }

  String next() {
    return Intl.message(
      'Next',
      name: 'next',
      desc: 'Next',
      locale: locale.toString(),
    );
  }

  String supportRole() {
    return Intl.message(
      'Support Role',
      name: 'supportRole',
      desc: 'Support Role',
      locale: locale.toString(),
    );
  }

  String provideSupportRole() {
    return Intl.message(
      'If we cannot satisfy your chosen conditions, would you consider providing a support role to those willing to go on strike?',
      name: 'provideSupportRole',
      desc:
          'If we cannot satisfy your chosen conditions, would you consider providing a support role to those willing to go on strike?',
      locale: locale.toString(),
    );
  }

  String yes() {
    return Intl.message(
      'Yes',
      name: 'yes',
      desc: 'Yes',
      locale: locale.toString(),
    );
  }

  String no() {
    return Intl.message(
      'No',
      name: 'no',
      desc: 'No',
      locale: locale.toString(),
    );
  }

  // SUPPORT ROLES

  String supportRoles() {
    return Intl.message(
      'Support Roles',
      name: 'supportRoles',
      desc: 'Support Roles',
      locale: locale.toString(),
    );
  }

  String pleaseSelectAnOptionFromTheList() {
    return Intl.message(
      'Please select an option from the list',
      name: 'pleaseSelectAnOptionFromTheList',
      desc: 'Please select an option from the list',
      locale: locale.toString(),
    );
  }

  String pleaseClickToSelectAvailableOptions() {
    return Intl.message(
      'Please click to select available options',
      name: 'pleaseClickToSelectAvailableOptions',
      desc: 'Please click to select available options',
      locale: locale.toString(),
    );
  }

  String hours() {
    return Intl.message(
      'hours',
      name: 'hours',
      desc: 'hours',
      locale: locale.toString(),
    );
  }

  String minimumPledgedTime() {
    return Intl.message(
      'Minimum pledged time',
      name: 'minimumPledgedTime',
      desc: 'Minimum pledged time',
      locale: locale.toString(),
    );
  }

  String typeYourAnswerHere() {
    return Intl.message(
      'Tap your answer here',
      name: 'typeYourAnswerHere',
      desc: 'Tap your answer here',
      locale: locale.toString(),
    );
  }

  String userNeeds() {
    return Intl.message(
      'User Needs',
      name: 'userNeeds',
      desc: 'User Needs',
      locale: locale.toString(),
    );
  }

  String projectDashboard() {
    return Intl.message(
      'Project Dashboard',
      name: 'projectDashboard',
      desc: 'Project Dashboard',
      locale: locale.toString(),
    );
  }

  String copyLink() {
    return Intl.message(
      'Copy Link',
      name: 'copyLink',
      desc: 'Copy Link',
      locale: locale.toString(),
    );
  }
}
