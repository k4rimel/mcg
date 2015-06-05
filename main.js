window.onload = function() {
	var canvas = document.getElementById("canvas"),
		context = canvas.getContext("2d"),
		width = canvas.width = window.innerWidth,
		height = canvas.height = window.innerHeight,
		mousestop = null,
		savedPosition = null,
		ball = particle.create(width / 2, height / 2, 0, 0);
		ball.thrust = vector.create(0, 0);
		ball.angle = 0;
		ball.turningLeft = false;
		ball.turningRight = false;
		ball.thrusting = false;
		ball.radius = 30;
		ball.setMass(100);

		ball2 = particle.create(750, 600, 0, 0);
		ball2.thrust = vector.create(0, 0);
		ball2.angle = 0;
		ball2.turningLeft = false;
		ball2.turningRight = false;
		ball2.thrusting = false;
		ball2.radius = 10;
		ball2.setMass(1);


	var gamestate = {
		value : null,
		set:function(state) {
			this.value = state;
		},
		get:function() {
			if(this.value)
				return this.value;
		}
	};
	update();
	function getMousePos(canvas, evt) {
	    var rect = canvas.getBoundingClientRect();
	    return {
	        X: evt.clientX - rect.left,
	        Y: evt.clientY - rect.top
	    };
 	}
 	function isColliding(ball1, ball2) {
 		var colliding = false;
 		getDistanceBetweenTwoParticles(ball1, ball2);
 		if(getDistanceBetweenTwoParticles(ball1, ball2) <= ball1.radius + ball2.radius) {
 			colliding = true;
 		}
 		return colliding;
 	}
 	function addBall() {
 		ball = particle.create(width / 2, height / 2, 0, 0),
 		ball.radius = 10;
 		ball.thrust = vector.create(0, 0),
 		ball.angle = 0,
 		ball.turningLeft = false,
 		ball.turningRight = false,
 		ball.thrusting = false;
 	}
 	function getDistanceBetweenTwoParticles(particle1, particle2) {
		var distX =   particle1.position.getX() - particle2.position.getX(),  
        distY   =   particle1.position.getY() - particle2.position.getY();
        return Math.sqrt(distX*distX + distY*distY);
 	}
 	function followCursor(pos) {
		var distX =   pos.X - ball.position.getX(),  
        distY   =   pos.Y - ball.position.getY();
        var dist = Math.sqrt(distX*distX + distY*distY);
        if(dist > ball.radius) {
        	ball.thrusting = true;
        	ball.angle  =   Math.atan2(distY, distX);
        } else {
        	ball.thrusting = false;
        }
 	}
 	function saveMousePos(pos) {
 		
 	}
 	canvas.addEventListener('mousemove', function(evt) {
 		var mousePos = getMousePos(canvas, evt);
 		if(mousestop) {
 			clearTimeout(mousestop);
 		}
 		mousestop =  setTimeout(function(){
 			savedPosition = mousePos;
 			followCursor(savedPosition);
 		}, 300);
 		 			savedPosition = mousePos;

	    followCursor(mousePos);

 	}, false);
	document.body.addEventListener("keydown", function(event) {
		// console.log(event.keyCode);
		switch(event.keyCode) {
			case 38: // up
				ball.thrusting = true;
				break;
				
			case 37: // left
				ball.turningLeft = true;
				break;
				
			case 39: // right
				ball.turningRight = true;
				break;			

			case 13: // enter : debug
				break;

			default:
				break;
				
		}
	});

	document.body.addEventListener("keyup", function(event) {
		// console.log(event.keyCode);
		switch(event.keyCode) {
			case 38: // up
				ball.thrusting = false;
				break;
				
			case 37: // left
				ball.turningLeft = false;
				break;
				
			case 39: // right
				ball.turningRight = false;
				break;

			default:
				break;
				
		}
	});
	function drawBall(obj, stopColor, thrustColor) {
		context.save();
		context.beginPath();
		context.arc(obj.position.getX(), obj.position.getY(), obj.radius, 0, Math.PI * 2, true);
		obj.accelerate(obj.thrust);
		obj.update();
		if(obj.thrusting) {
			context.strokeStyle = thrustColor;
			context.fillStyle = thrustColor;
		} else {
			context.strokeStyle = stopColor;
			context.fillStyle = stopColor;
		}
		context.fill();
		context.stroke();
		context.restore();
		context.closePath();
	}
	function drawSpeed() {
		if(ball.thrusting) {
			context.fillStyle = "red";
		} else {

			context.fillStyle = "blue";
		}
	  	context.font = "bold 16px Arial";
	  	context.fillText('Speed : '+ (ball.velocity.getLength()).toFixed(3), 100, 100);
	}
	function drawArena(radius) {
		context.beginPath();
      	context.arc(width/2, height/2, radius, 0, 2 * Math.PI, false);
      	context.fillStyle = '#4679BD';
      	context.fill();
      	// context.lineWidth = 5;
      	context.strokeStyle = '#4679BD';
      	context.stroke();
      	context.closePath();
	}
	function getNVector(ball1, ball2) {
		nVector = vector.create(ball2.position.getX() - ball1.position.getX(),
								ball2.position.getY() - ball1.position.getY());
		return nVector;
	}
	// function getUNVector(nVector) {
	// 	unVector = nVector.divide(Math.sqrt(Math.pow(nVector.getX(), 2), Math.pow(nVector.getY(), 2)));
	// 	return unVector;
	// }
	function getUTVector(unVector) {
		utVector = vector.create(-1*unVector.getY(), unVector.getX());
		return utVector;
	}
	function collide(ball1, ball2) {

		// unVector = getNVector(ball1, ball2);
		// // unVector = getUNVector(getNVector(ball1, ball2));
		// utVector = getUTVector(getNVector(ball1, ball2));

		// v1 = ball1.velocity;
		// v2 = ball2.velocity;

		// m1 = ball1.mass;
		// m2 = ball2.mass;

		// v1n = unVector.dotProduct(v1);
		// v2n = unVector.dotProduct(v2);

		// v1t = utVector.dotProduct(v1);
		// v2t = utVector.dotProduct(v2);

		// v1nPrime = ( ( v1n*(m1-m2) ) + (2*m2*v2n) ) / (m1+m2);
		// v2nPrime = ( ( v2n*(m2-m1) ) + (2*m1*v1n) ) / (m1+m2);
		// v1tPrime = v1t;
		// v2tPrime = v2t;

		// v1nPrimeV = unVector.multiply(v1nPrime);
		// v2nPrimeV = unVector.multiply(v2nPrime);
		// v1tPrimeV = utVector.multiply(v1tPrime);
		// v2tPrimeV = utVector.multiply(v2tPrime);

		// v1Prime = v1nPrimeV.addTo(v1tPrimeV);
		// v2Prime = v2nPrimeV.addTo(v2tPrimeV);
		// if(v1Prime !== 'undefined' && v2Prime !== 'undefined') {
		// 	console.log((2*m2*v2n));



		// 	console.log(unVector);
		// 	console.log(v1nPrime);
		// 	console.log(v1tPrimeV);
		// 	console.log(v1Prime);
		// } else {
		// 	console.log(v1nPrimeV);
		// 	console.log(v1tPrimeV);
		// 	console.log(v2nPrimeV);
		// 	console.log(v2tPrimeV);
		// }

		// ball1.velocity = v1Prime;
		// ball2.velocity = v2Prime;
		var xDistance = (ball2.position.getX() - ball1.position.getX());
		var yDistance = (ball2.position.getY() - ball1.position.getY());

		var normalVector = vector.create(xDistance, yDistance); // normalise this vector store the return value in normal vector.
		normalVector = normalVector.normalise();

		var tangentVector = vector.create((normalVector.getY() * -1), normalVector.getX());
		
		// create ball scalar normal direction.
		var ball1scalarNormal =  normalVector.dot(ball1.velocity);
		var ball2scalarNormal = normalVector.dot(ball2.velocity);

		// create scalar velocity in the tagential direction.
		var ball1scalarTangential = tangentVector.dot(ball1.velocity); 
		var ball2scalarTangential = tangentVector.dot(ball2.velocity); 

		var ball1ScalarNormalAfter = (ball1scalarNormal * (ball1.getMass() - ball2.getMass()) + 2 * ball2.getMass() * ball2scalarNormal) / (ball1.getMass() + ball2.getMass());
		var ball2ScalarNormalAfter = (ball2scalarNormal * (ball2.getMass() - ball1.getMass()) + 2 * ball1.getMass() * ball1scalarNormal) / (ball1.getMass() + ball2.getMass());

		var ball1scalarNormalAfter_vector = normalVector.multiply(ball1ScalarNormalAfter); // ball1Scalar normal doesnt have multiply not a vector.
		var ball2scalarNormalAfter_vector = normalVector.multiply(ball2ScalarNormalAfter);

		var ball1ScalarNormalVector = (tangentVector.multiply(ball1scalarTangential));
		var ball2ScalarNormalVector = (tangentVector.multiply(ball2scalarTangential));;

		ball1.velocity = ball1ScalarNormalVector.add(ball1scalarNormalAfter_vector);
		ball2.velocity = ball2ScalarNormalVector.add(ball2scalarNormalAfter_vector);

		// ball1.position = ball1.lastGoodPosition;
		// ball2.position = ball2.lastGoodPosition;
	}
	function update() {

		context.clearRect(0, 0, width, height);
	  	drawArena(width/4);

		drawBall(ball, "#AAAAAA", "#FF6961");
		ball.thrust.setAngle(ball.angle);
		if(ball.thrusting) {
			ball.thrust.setLength(0.1);
		}
		else {
			ball.thrust.setLength(0);
		}

		drawBall(ball2, "#AAAAAA", "#77DD77");
		ball2.thrust.setAngle(ball2.angle);
		if(ball2.thrusting) {
			ball.thrust.setLength(0.1);
		}
		else {
			ball2.thrust.setLength(0);
		}

		drawSpeed();
		if(mousestop) {
			if(savedPosition) {
				followCursor(savedPosition);
			}
		}
		if(ball.position.getX() > width || ball.position.getX() < 0 || ball.position.getY() > height || ball.position.getY() < 0) {
			gamestate.set('end');
		}
		if(isColliding(ball, ball2)) {
			console.log('collision');
			collide(ball, ball2);
		}
		// requestAnimationFrame();
	}
	setInterval(update, 10);
};