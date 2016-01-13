// Package bruteblocker allows to slow down bruteforce attacks by delaying
// request responses or blocking further authentication attempts for some
// time.
// This makes it harder for the attacker to guess passwords or user names,
// while retaining an acceptable amount of usability for normal users.
//
// The idea is to keep track of the number of failed authentication attempts
// as well as the timestamp of the last authentication attempt.
// This is done per user, per IP address and globally.
//
// Per failed authentication attempt, the delay is increased.
// After a while, the delay is reset to zero.
package bruteblocker
