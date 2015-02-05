var createJSErrorUrl = function (message, url, line_num, column) {
    var base_url = 'http://127.0.0.1:8080/jsec?r=';
    var err_obj = {
        m: message,
        u: url,
        l: line_num,
        c: column,
    };
    var err_json = JSON.stringify(err_obj);
    var err_string = encodeURIComponent(err_json);
    return base_url + err_string;
};

window.onerror = function (message, url, line_num, column) {
  (new Image()).src = createJSErrorUrl(message, url, line_num, column);
};

throw new Error('test');
