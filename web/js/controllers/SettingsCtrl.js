'use strict';

/* Controllers */

angular.module('xtv.controllers').
  controller('SettingsCtrl', ['$scope', '$http', function($scope, $http) {

    $scope.loadSettings = function () {
      $http.get('data/loadSettings').success(function (response) {
        if (response.status == 'ok') {
          $scope.settings = response.result;

        } else {
          msg.error(response.message);
        }
      });
    };
    $scope.loadSettings();

    $scope.saveSettings = function () {
      $http.post('data/saveSettings', {data: $scope.settings}).success(function (response) {
        if (response.status = 'ok') {
          $scope.hideSettingsDialog();
          $scope.loadSettings();
        } else {
          msg.error(response.message);
        }
      });
    };

    $scope.$on('xtv:showSettingsDialog', function () {
      $scope.showSettingsDialog();
    });

    $scope.showSettingsDialog = function () {
      $('#settingsDialog').modal('show');
    };

    $scope.hideSettingsDialog = function () {
      $('#settingsDialog').modal('hide');
    };

  }]);