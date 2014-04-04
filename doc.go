// Package goj is a tool to filter decoded json.
//
// goj attempts to match the jsonpath syntax
//
// 	JSONPath         | goj              | Description
// 	---------------- | ---------------- | -----------
// 	$                | n/a              | The root of the path is implied. Optionally can start with path separator.
// 	. or []          | . or []          | child operator
// 	n/a              | ..               | parent operator
// 	*                | *                | wildcard.
// 	n/a              | **               | recursive wildcard.
// 	[]               | []               | native array operator
// 	[,]              | [,]              | array indices as a set. (not sure if JSONPath supports this for names.)
// 	[start:stop:end] | [start:stop:end] | array slice operator.
// 	?() or ()        | n/a              | script expression is not supported yet.
//
// The major difference is that JSONPath is primarly used to retrieve values, while goj is used to *filter* values.
//
// Given a json structure like this:
//   {
//     "store": {
//       "book": [
//         {
//           "category": "reference",
//           "author": "Nigel Rees",
//           "title": "Sayings of the Century",
//           "price": 8.95
//         },
//         {
//           "category": "fiction",
//           "author": "Evelyn Waugh",
//           "title": "Sword of Honour",
//           "price": 12.99
//         }
//       ]
//     }
//   }
//
// A JSONPath like this:
// 	$.store.book[1].author
//
// would return a *value* of "Evelyn Waugh"
//
// A similar goj path of:
// 	store.book[1].author
//
// would return a filterd structure:
//   {
//     "store": {
//       "book": [
//         {
//           "author": "Evelyn Waugh"
//         }
//       ]
//     }
//   }
//
// Two major additions to the syntax are the parent operator `..` and the recursive wildcard `**`.
//
// A note about the path separtor. I chose '.' instead of '/' or something else because I wanted to follow JSONPath syntax.
package goj
