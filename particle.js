var particle = {
	position: null,
	velocity: null,
	friction: null,
	mass: 1,

	create: function(x, y, speed, direction) {
		var obj = Object.create(this);
		obj.position = vector.create(x, y);
		obj.velocity = vector.create(0, 0);
		obj.friction = 0.970;
		obj.velocity.setLength(speed);
		obj.velocity.setAngle(direction);
		return obj;
	},

	accelerate: function(accel) {
		this.velocity.addTo(accel);
	},	
	applyFriction: function() {
		if(this.velocity.getLength() < 0.250) {
			this.velocity.setLength = 0;
		}
		this.velocity.multiplyBy(this.friction);
	},
	destroy: function() {
		console.log('lost');
	},
	setMass: function (inMass) { 
		this.mass = inMass;
	},
	getMass: function () { 
		return this.mass;
	},
	update: function() {
		this.position.addTo(this.velocity);
		this.applyFriction();
	},

};