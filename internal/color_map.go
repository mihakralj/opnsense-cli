/*
Copyright © 2023 Miha miha.kralj@outlook.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package internal

var c = map[string]string{
	"tag": "\033[0m",
	"txt": "\033[33m",
	"atr": "\033[33m",

	"chg": "\033[36m",
	"add": "\033[32m",
	"del": "\033[31m\033[9m",
	"red": "\033[31m",
	"grn": "\033[32m",
	"ele": "\033[36m",

	"yel": "\033[33m",
	"blu": "\033[34m",
	"mgn": "\033[35m",
	"cyn": "\033[36m",
	"wht": "\033[37m",
	"gry": "\033[90m",

	"ita": "\033[3m", // italics
	"bld": "\033[1m", // bold
	"stk": "\033[9m", // strikethroough
	"und": "\033[4m",
	"rev": "\033[7m", // reverse colors

	"ell": "\u2026",
	"arw": " \u2192 ",
	"nil": "\033[0m",
}
