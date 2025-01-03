// memguardian is a package that provides a simple RAM memory control mechanism
// once activated it sets an internal atomic boolean when the RAM usage exceed in absolute
// terms the warning ratio, for passive indirect check or invoke an optional callback for
// reactive backpressure
package memguardian
