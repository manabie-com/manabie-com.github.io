"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.filterEqual = exports.filterExisting = exports.parseFile = exports.readFile = void 0;
const fs = __importStar(require("fs"));
const mime = __importStar(require("mime"));
const YAML = __importStar(require("js-yaml"));
const readFile = (file) => new Promise((resolve, reject) => fs.readFile(file, (err, v) => err ? reject(err) : resolve(v.toString('utf-8'))));
exports.readFile = readFile;
const parseFile = (file, data) => mime.getType(file) === 'text/yaml' ? YAML.safeLoad(data) : JSON.parse(data);
exports.parseFile = parseFile;
const filterExisting = (keepExisting, existingArray = []) => (item) => {
    if (keepExisting) {
        const found = existingArray.find((t) => t.name === item.name);
        return !found;
    }
    return true;
};
exports.filterExisting = filterExisting;
const filterEqual = (existingArray = []) => (item) => {
    const toggle = existingArray.find((t) => t.name === item.name);
    if (toggle) {
        return JSON.stringify(toggle) !== JSON.stringify(item);
    }
    return true;
};
exports.filterEqual = filterEqual;