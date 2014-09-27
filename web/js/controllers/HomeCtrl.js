'use strict';

/* Controllers */

angular.module('xtv.controllers').
  controller('HomeCtrl', ['$scope', '$http', 'msg', function($scope, $http, msg) {

    $scope.loadServers = function () {
      $http.get('data/loadServers').success(function(response){
        if (response.status == 'ok') {
          $scope.servers = response.results;

          //reselect server
          angular.forEach ($scope.servers, function (server) {
             $scope.getServerStatus (server);
             if ($scope.selectedServer) {
               if (server.id == $scope.selectedServer.id) {
                 $scope.selectedServer = server;
                 $scope.loadConsole (server);
               }
             }
          });
        } else {
          msg.error (response.message);
        }
      });
    };
    $scope.loadServers();

    $scope.selectServer = function (server) {
      $scope.selectedServer = server;
      $scope.loadConsole (server);
    };

    $scope.showAddServerDialog = function () {
      $('#addServerDialog').modal('show');
    };

    $scope.addServer = function () {
      var newServer = {
          name: $scope.newServer.uri,
          port: $scope.newServer.port,
          status: 'Not Connected',
          channels: []
      };

      $http.post ('data/saveServer', {data: newServer}).success (function (response) {
        if (response.status = 'ok') {
          $('#addServerDialog').modal('hide');
          $scope.newServer = undefined;
          $scope.loadServers();
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.showAddChannelDialog = function () {
      $('#addChannelDialog').modal('show');
    };

    $scope.addChannel = function () {
      var data = {
        serverId: $scope.selectedServer.id,
        channel: {
          name: $scope.newChannel.name
        }
      };
      $http.post ('data/addChannel', {data: data}).success (function (response) {
        if (response.status = 'ok') {
          $('#addChannelDialog').modal('hide');
          $scope.newChannel = undefined;
          $scope.loadServers();
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.showDeleteServerConfirm = function (server) {
      $scope.serverToDelete = server;
      $('#deleteServerConfirmDialog').modal ('show');
    };

    $scope.deleteServer = function () {
      $http.post ('data/deleteServer', {data: $scope.serverToDelete}).success (function (response) {
        if (response.status = 'ok') {
          $scope.selectedServer = undefined;
          $scope.serverToDelete = undefined;
          $scope.selectedServerConsole = undefined;
          $scope.loadServers();
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.showDeleteChannelConfirm = function (channel) {
      $scope.channelToDelete = channel;
      $('#deleteChannelConfirmDialog').modal ('show');
    };

    $scope.deleteChannel = function () {
      var data = {
        serverId: $scope.selectedServer.id,
        channelId: $scope.channelToDelete.id
      };

      $http.post ('data/deleteChannel', {data: data}).success (function (response) {
        if (response.status = 'ok') {
          $scope.channelToDelete = undefined;
          $scope.loadServers();
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.getStatusClass = function (server) {
      if (server.status == 'Connected') {
         return 'glyphicon-globe';
      }
      return 'glyphicon-ban-circle';
    };

    $scope.toggleConnection = function (server) {
      $http.post ('irc/toggleConnection', {data: angular.copy (server)}).success (function (response) {
        if (response.status = 'ok') {
          server.status = response.result.status;
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.getServerStatus = function (server) {
      $http.post ('irc/getServerStatus', {data: angular.copy (server)}).success (function (response) {
        server.status = response.status;
      });
    };

    $scope.loadConsole = function (server) {
      $scope.selectedServerConsole = undefined;
      $http.post ('irc/getServerConsole', {data: angular.copy (server)}).success (function (response) {
        $scope.selectedServerConsole = response.console;
      });
    };
  }]);
