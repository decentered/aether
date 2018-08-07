"use strict";
// Globals > Methods
// This file contains methods that are useful in multiple places, and depend on nothing but themselves.
module.exports = {
    GetUserName: function (owner) {
        if (typeof owner === 'undefined') {
            return "";
        }
        if (typeof owner.fingerprint === 'undefined') {
            // Necessary because the 'observer' object is not undefined, but also not what we want.
            return "";
        }
        if (owner.compiledusersignals.cnamesourcefingerprint.length > 0 && owner.compiledusersignals.canonicalname.length > 0) {
            return owner.compiledusersignals.canonicalname;
        }
        if (owner.noncanonicalname.length > 0) {
            return owner.noncanonicalname;
        }
        return owner.fingerprint;
    },
    TimeSince: function (timestamp) {
        var now = new Date();
        var ts = new Date(timestamp * 1000);
        var secondsPast = (now.getTime() - ts.getTime()) / 1000;
        if (secondsPast < 60) {
            return secondsPast + 's';
        }
        if (secondsPast < 3600) {
            return (secondsPast / 60) + 'm';
        }
        if (secondsPast <= 86400) {
            return (secondsPast / 3600) + 'h';
        }
        // If older than a day
        var day;
        var month;
        var year;
        day = ts.getDate();
        var tsds = ts.toDateString().match(/ [a-zA-Z]*/);
        if (tsds !== null) {
            month = tsds[0].replace(" ", "");
        }
        year = ts.getFullYear() === now.getFullYear() ? "" : " " + ts.getFullYear();
        return day + " " + month + year;
    }
};
//# sourceMappingURL=methods.js.map