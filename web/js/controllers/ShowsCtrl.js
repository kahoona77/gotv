'use strict';

/* Controllers */

angular.module('xtv.controllers').
  controller('ShowsCtrl', ['$scope', 'msg', '$http', '$location', function($scope, msg, $http, $location) {

    $scope.loadShows = function () {
      $http.get('/shows/load').success(function(response){
        if (response.status == 'ok') {
          $scope.shows = response.shows;

          //reselect server
          angular.forEach ($scope.shows, function (show) {
            if ($scope.selectedShow) {
              if (show._id == $scope.selectedShow._id) {
                $scope.selectedShow = show;
              }
            }
          });
        } else {
          msg.error (response.message);
        }
      });
    };
    $scope.loadShows();

    $scope.showAddShowDialog = function () {
      $('#addShowDialog').modal('show');
    };

    $scope.selectShow = function (show) {
      $scope.selectedShow = show;
      $scope.loadEpisodes (show);
    };

    $scope.loadEpisodes = function (show) {
     $scope.seasons = undefined;
     $http.get('/shows/loadEpisodes', {params : {showId: show._id}}).success(function(response){
        if (response.status == 'ok') {
          $scope.seasons = response.data;
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.showDeleteShowConfirm = function (show) {
      $scope.showToDelete = show;
      $('#deleteShowConfirmDialog').modal ('show');
    };

    $scope.deleteShow = function () {
      $http.post ('shows/delete', {data: $scope.showToDelete}).success (function (response) {
        if (response.status = 'ok') {
          $('#deleteShowConfirmDialog').modal('hide');
          $scope.showToDelete = undefined;
          $scope.selectedShow = undefined;
          $scope.seasons = undefined;
          $scope.loadShows();
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.searchEpisode = function (show, episode) {
      var pad = "00";
      var season = "" + episode.season;
      season =  pad.substring(0, pad.length - season.length) + season;

      var number = "" + episode.number;
      number =  pad.substring(0, pad.length - number.length) + number;

      var query = show.name + " S" +  season + "E" + number;
      $location.url ('/search/' + query);
    };

    //search
    $scope.searchShow = function () {

      $http.get('/shows/search', {params : {query: $scope.query}}).success(function(response){
        if (response.status == 'ok') {
          $scope.searchResults = response.shows;
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.addShow = function (show) {
      $http.post ('/shows/add', {data: show}).success (function (response) {
        if (response.status = 'ok') {
          $('#addShowDialog').modal('hide');
          $scope.query = undefined;
          $scope.loadShows();
        } else {
          msg.error (response.message);
        }
      });
    };

    $scope.markEpisode = function (episode) {
      episode.loaded = !episode.loaded;
      $http.post ('/shows/saveEpisode', {data: episode}).success (function (response) {
        if (response.status != 'ok') {
          msg.error (response.message);
        }
      });
    };

    //Filter
    $scope.hideLoaded = true;
    $scope.showLoaded = function () {
      $scope.hideLoaded = !$scope.hideLoaded;
    };

    $scope.showLoadedEpisode = function(episode) {
      if ($scope.hideLoaded && episode.loaded) {
        return false;
      }
      return true;
    };

    $scope.episodesLeftToShow = function (episodes) {
      for (var i = 0; i < episodes.length; i++) {
        var ep = episodes[i];
        if ($scope.showLoadedEpisode (ep)) {
          return true;
        }
      }
      return false;
    };

    $scope.markShowAsLoaded = function (show) {
      $http.post ('/shows/markShowLoaded', {showId: show._id}).success (function (response) {
        if (response.status != 'ok') {
          msg.error (response.message);
        } else {
          $scope.loadEpisodes ($scope.selectedShow);
        }
      });
    };

    $scope.markSeasonAsLoaded = function (show, season) {
      $http.post ('/shows/markSeasonLoaded', {showId: show._id, season: season}).success (function (response) {
        if (response.status != 'ok') {
          msg.error (response.message);
        } else {
          $scope.loadEpisodes ($scope.selectedShow);
        }
      });
    };


  }]);
