if (typeof window === 'undefined') {
  var Config = require('./conf.js'),
    Lib = require('./lib.js'),
    User = require('./user.js').User,
    MsgUser = require('./user.js').MsgUser
}


function Users(opts) {
  opts = opts || {};

  this.users = {};
  this.usersList = [];
  this.userMap = {};  // map references in obj to list

  var self = this;

  // returns if a user is stored
  this.hasUser = function(userId) {
    return userId in self.users;
  };

  // returns a user given a userid, null if the user dne
  this.getUser = function(userId) {
    if (userId in self.users) {
      return self.users[userId];
    }
    return null;
  };

  // returns the number of users stored
  this.numUsers = function() {
    return self.users.length;
  };

  this.addUser = function(user) {
    if (user.id in self.users) {
      console.warn('user exists, overwriting');
    }

    self.users[user.id] = user;
    self.userMap[user.id] = self.usersList.length;
    self.usersList.push(user);
  };

  this.removeUser = function(userId) {
    if (self.hasUser(userId)) {
      delete self.users[userId];
    }
    else {
      console.warn('remove: user dne');
    }
  };

  this.addFromMsgUser = function(msgUser) {
    var user = User.fromMsgUser(msgUser);
    self.addUser(user);
  };

  // if a user with the same id as msgUser then update it,
  // else create a new use from msgUser
  this.updateFromMsgUser = function(msgUser) {
    if (msgUser.isInvalid()) {
      console.warn('msgUser invalid attrs: ', msgUser.isInvalid());
    }
    console.assert(!msgUser.isInvalid());

    var user = self.getUser(msgUser);
    if (user) {
      user.updateFromMsgUser(msgUser);
    }
    else {
      self.addFromMsgUser(msgUser);
    }
    return true;
  };

  // updates a stored user given a msgUser
  this.updateUser = function(msgUser) {
    var user = self.getUser(msgUser.id);
  };
};


if (typeof window === 'undefined') {
  module.exports = Users;
}
