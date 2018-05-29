<?php
/* Snap4city IOT Application API
   Copyright (C) 2018 DISIT Lab http://www.disit.org - University of Florence

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as
   published by the Free Software Foundation, either version 3 of the
   License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>. */

function new_nodered($db,$uname, $aname, $image) {
  include '../config.php';
  
  $aid=random_str($app_id_length);          
  $name = "nr".$aid;

  //prepare the /data dir for nodered
  $out = array();
  $ret = null;
  exec($nodered_script.' '.$name.' '.$uname,$out,$ret);

  $result=new_marathon_nodered_container($marathon_url,$image,$name,$uname);
  if($result) {
    $did = "";
    if(is_string($result)) {
      $did = $result;
    }
    else if(count($result["tasks"])==0){
      $did = $result["deployments"][0]["id"];
    } else {
      $id = $result["tasks"][0]["id"];
    }
    
    $proxy='        location  /nodered/'.$name.'/ {
                proxy_set_header Host $http_host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection "upgrade";
                proxy_pass "http://iot-app.snap4city.org/nodered/'.$name.'/";
        }';
    $f=fopen("/mnt/data/proxy/nr/$name.conf","w");
    if($f) {
      fwrite($f,$proxy);
      fclose($f);
    }
    
    $url = "https://iot-app.snap4city.org/nodered/$name";
    return array("id"=>$name,"url"=>$url,'name'=>$aname);
  } else {
    //mysqli_query($db, "UPDATE application SET status='ERROR',status_description='cannot create container' WHERE id=$aid") or die(mysqli_error($db));
    return array('error'=>"cannot create container");
  }
}

function new_plumber($db,$uid) {
  /*
  if($r=mysqli_query($db, "SELECT COUNT(*) as napps FROM application WHERE uid='$uid' AND (status='RUNNING' OR status='STOPPED')") or die(mysqli_error($db))) {
    if($o=mysqli_fetch_object($r)) {
      if($o->napps >= 10) {
        echo "{error:\"cannot create you reached maximum number\"}";
        return;
      }
    }
  }
  
  mysqli_query($db, "INSERT INTO application(uid,status) VALUES($uid,'PRE-CREATE')");
  $aid=mysqli_insert_id($db);
  $name = "plumber$aid";

  $result=  new_marathon_plumber_container("192.168.1.187",null,$name);
  if($result) {
    $did = "";
    if(is_string($result)) {
      mysqli_query($db, "UPDATE application SET status='DEPLOYING', name='$name',deployment_id='$result' WHERE id=$aid") or die(mysqli_error($db));
    }
    else if(count($result["tasks"])==0){
      $did = $result["deployments"][0]["id"];
      mysqli_query($db, "UPDATE application SET status='DEPLOYING', name='$name',deployment_id='$did' WHERE id=$aid") or die(mysqli_error($db));
    } else {
      var_dump($result);
      $id = $result["tasks"][0]["id"];
      mysqli_query($db, "UPDATE application SET status='RUNNING', name='$name',container_id='$id' WHERE id=$aid") or die(mysqli_error($db));
    }
    $url = "http://www.snap4city.org/plumber/$name";
    echo "{\"id\":\"$aid\",\"url\":\"$url\"}";
  } else {
    mysqli_query($db, "UPDATE application SET status='ERROR',status_description='cannot create container' WHERE id=$aid") or die(mysqli_error($db));
    echo "{error:\"cannot create container\"}";
  } 
   *
   */ 
}

function stop_nodered($db,$uid,$aid) {
  /*
  if($r=mysqli_query($db, "SELECT container_id,host FROM application WHERE id='$aid' AND uid='$uid' AND status='RUNNING'") or die(mysqli_error($db))) {
    if($o=mysqli_fetch_object($r)) {
    } else {
        echo "{error:\"cannot find running container $aid for user $uid\"}";      
    }
  }*/
}

function start_nodered($db,$uid,$aid) {
  /*
  if($r=mysqli_query($db, "SELECT container_id,host FROM application WHERE id='$aid' AND uid='$uid' AND status='STOPPED'") or die(mysqli_error($db))) {
    if($o=mysqli_fetch_object($r)) {
      $host = $o->host;
      $id= $o->container_id;
      $result = http_post("http://$host:3000/v1.30/containers/$id/start","","");
      if($result["httpcode"]==204) {
        mysqli_query($db, "UPDATE application SET status='RUNNING',status_description='' WHERE id=$aid") or die(mysqli_error($db));
        echo "{}";
      } else {
        $msg = $result["result"]["message"];
        echo "{error:\"cannot stop container $msg\"}";
      }
    } else {
        echo "{error:\"cannot find stopped container $aid for user $uid\"}";      
    }
  } 
   * 
   */ 
}

function rm_app($db,$uid,$aid) {
  include '../config.php';
  
  $result = http_delete($marathon_url."/v2/apps/".$aid);
  if($result["httpcode"]==200) {
    echo json_encode($result["result"]);
  }
  else {
    echo "{error:\"failed removal of app $aid\"}";
  }  
}

function restart_app($db,$uid,$aid) {
  include '../config.php';
  $result = http_post($marathon_url."/v2/apps/".$aid."/restart","","application/json");
  if($result["httpcode"]==200) {
    echo json_encode($result["result"]);
  }
  else {
    echo "{error:\"failed restart of app $aid\"}";
  }  
}

function status_app($db,$uid,$aid) {
  include '../config.php';
  $result = http_get($marathon_url."/v2/apps/".$aid);
  if($result["httpcode"]==200) {
    $r = array();
    @$r['healthiness']=$result['result']['app']['tasks'][0]['healthCheckResults'][0]['alive'] ?: false;
    echo json_encode($r);
  }
  else {
    echo "{error:\"failed access to status for app $aid\"}";
  }  
}

function new_marathon_nodered_container($base_url,$image, $id, $uname) {
  include '../config.php';
//    "cmd": "/bin/bash -c \"cd /usr/src/node-red ; npm start -- --userDir /data\"",

  $json = '{
    "id": "'. $id .'",
    "cmd": "npm start -- --userDir /data",
    "labels": {
      "HAPROXY_GROUP":"external",
      "HAPROXY_0_VHOST":"iot-app.snap4city.org",
      "HAPROXY_0_PATH":"/nodered/'. $id .'/"
    },
    "container": {
        "type": "DOCKER", 
        "docker": {
            "network": "HOST", 
            "image": "'. $image .'"
        },
        "volumes": [
        {
            "containerPath": "/data",
            "hostPath": "/mnt/data/nr-data/'. $id .'",
            "mode": "RW"
        }
    ]        
    },
    "cpus": '.$nodered_cpu.', 
    "portDefinitions": [{"name": null, "protocol": "tcp", "port": 0, "labels": null}],
    "instances": 1,
    "constraints": [["@hostname", "UNLIKE", "mesos[1-3]t"]],
    "env": {}, 
    "mem": '.$nodered_mem.', 
    "disk": 128, 
    "healthChecks": [{
        "maxConsecutiveFailures": 0, 
        "protocol": "HTTP", 
        "portIndex": 0, 
        "gracePeriodSeconds": 240, 
        "path": "/nodered/'.$id.'/ui", 
        "timeoutSeconds": 10, 
        "intervalSeconds": 15}], 
    "appId": "'. $id .'"
}';
  /*$out = array();
  $ret = null;
  exec("/home/ubuntu/add-nodered.sh $id $uname",$out,$ret);*/
  $result = http_post($base_url."/v2/apps",$json, "application/json");
  if($result["httpcode"]==201) {
    store_on_disces_em($id, $json);
    return $result["result"];
  }
  else if($result["httpcode"]==200) {
    store_on_disces_em($id, $json);
    return $result["result"]["deploymentId"];
  } else {
    //var_dump($result);
    return "";
  }
  return "";
}

function new_marathon_plumber_container($base_url,$image, $id) {
  $json = '{
    "cmd": "Rscript /root/Snap4City/Snap4CityStatistics/RunRestApi.R",
    "id": "'. $id .'",
    "labels": {
      "HAPROXY_GROUP":"external",
      "HAPROXY_0_VHOST":"nr.snap4city.org",
      "HAPROXY_0_PATH":"/plumber/'. $id .'/",
      "HAPROXY_0_HTTP_BACKEND_PROXYPASS_PATH":"/plumber/'. $id .'/"
    },
    "container": {
        "type": "DOCKER", 
        "docker": {
            "network": "BRIDGE", 
            "image": "disit/plumber:version7",
            "portMappings": [
              {
                "containerPort": 8080,
                "hostPort": 0,
                "protocol": "tcp",
                "name": "http"
              }]
        },
        "volumes": [
        {
            "containerPath": "/root",
            "hostPath": "/mnt/R",
            "mode": "RW"
        }
    ]        
    },
    "cpus": 0.095, 
    "portDefinitions": [{"name": null, "protocol": "tcp", "port": 0, "labels": null}],
    "instances": 1, 
    "env": {}, 
    "mem": 140, 
    "disk": 128, 
    "healthChecks": [{
        "maxConsecutiveFailures": 0, 
        "protocol": "HTTP", 
        "portIndex": 0, 
        "gracePeriodSeconds": 240, 
        "path": "/sum?a=1&b=1", 
        "timeoutSeconds": 10, 
        "intervalSeconds": 15}], 
    "appId": "'. $id .'"
}';
  $out = array();
  $ret = null;
  //exec("/home/ubuntu/add-nodered.sh $id",$out,$ret);
  $result = http_post($base_url."/v2/apps",$json, "application/json");
  if($result["httpcode"]==201)
    return $result["result"];
  else if($result["httpcode"]==200)
    return $result["result"]["deploymentId"];
  else {
    var_dump($result);
    return "";
  }
  return "";
}

function random_str($length, $keyspace = '0123456789abcdefghijklmnopqrstuvwxyz')
{
    $pieces = [];
    $max = strlen($keyspace) - 1;
    for ($i = 0; $i < $length; ++$i) {
        $pieces []= $keyspace[random_int(0, $max)];
    }
    return implode('', $pieces);
}
