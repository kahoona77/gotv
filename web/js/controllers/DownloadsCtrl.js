'use strict';

/* Controllers */

angular.module('xtv.controllers').
  controller('DownloadsCtrl', ['$scope', '$timeout', '$http', function($scope, $timeout, $http) {

    $scope.downloads = [];

    $scope.loadDownloads = function () {
        $http.get('downloads/listDownloads',{params: { 'nocache': new Date().getTime() }}).success(function(response){
            if (response.status == 'ok') {
                $scope.downloads = response.downloads;
            } else {
                msg.error (response.message);
            }
        });
    };
    $scope.loadDownloads();

    $scope.startReloadTimer = function () {
      $scope.reloadTimer = $timeout (function () {
          $scope.loadDownloads();
          $scope.startReloadTimer();
      }, 1000)
    };
    $scope.startReloadTimer();

    $scope.$on ('$locationChangeStart', function () {
      if ($scope.reloadTimer) {
        $timeout.cancel ($scope.reloadTimer);
      }
    });

    $scope.selectItem = function (item) {
      $scope.selectedItem = item;
    };

    $scope.stop = function () {
        $http.post('downloads/stopDownload', {data: $scope.selectedItem}).success(function(response){
            if (response.status == 'ok') {
                $scope.selectedItem = undefined;
                $scope.loadDownloads();
            } else {
                msg.error (response.message);
            }
        });
    };

    $scope.resume = function () {
        $http.post('downloads/resumeDownload', {data: $scope.selectedItem}).success(function(response){
            if (response.status == 'ok') {
                $scope.selectedItem = undefined;
                $scope.loadDownloads();
            } else {
                msg.error (response.message);
            }
        });
    };

    $scope.showCancelConfirm = function () {
      $('#confirmDialog').modal ('show');
    };

    $scope.cancel = function () {
        $http.post('downloads/cancelDownload', {data: $scope.selectedItem}).success(function(response){
            if (response.status == 'ok') {
                $scope.selectedItem = undefined;
                $scope.loadDownloads();
            } else {
                msg.error (response.message);
            }
        });
    };

    $scope.clear = function () {
      var completed = [];
      angular.forEach ($scope.downloads, function (item) {
         if (item.status == 'COMPLETE') {
           completed.push (item);
         }
      });

      angular.forEach (completed, function (item) {
         $http.post('downloads/cancelDownload', {data: item}).success(function(response){
              if (response.status == 'ok') {
                  $scope.selectedItem = undefined;
                  $scope.loadDownloads();
              } else {
                  msg.error (response.message);
              }
          });
      });
    };

    $scope.calcTimeRemaining = function (item) {
      var remainingKBytes = (item.size - item.bytesReceived) / 1024;
      var remainingSeconds = remainingKBytes / item.speed;

      var min = Math.floor(remainingSeconds / 60);
      var sec = Math.round(remainingSeconds % 60);

      if (min < 10) {
        min = '0' + min
      }

      if (sec < 10) {
        sec = '0' + sec
      }

      return min + ":" + sec + " Minutes";
    };
  }]);
