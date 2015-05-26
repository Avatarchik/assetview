app = angular.module('assetView', ["ngAnimate"]);
app.controller 'AssetListCtrl', ($scope, $http)->
  $scope.assets = []
  $scope.loading = true
  $http.get "/$"
    .success (r)=>
      $scope.assets = r
      $scope.loading = false
    .error ()=>
      console.log arguments

  $http.get "/@"
    .success (r)=>
      $scope.status = r
      console.log $scope.status
    .error ()=>
      console.log arguments

  $scope.$watch "detail", (c)->
    return unless c?.path?
    $http.get "/*/"+c.path
      .success (r)=>
        $scope.detail2 = r
      .error ()=>
        console.log arguments
    #http://localhost:12345/*/jp_v2/area/area_10.Android.unity3d
