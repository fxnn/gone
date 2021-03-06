
var GoneEditor = new(function() {

	var _root = this;
    var _elements = {
        content: null,
        editor: null,
		contentType: null,
        form: null
    };
	var _modesPerContentType = {
		"javascript": "javascript",
		"text/html": "html",
		"text/css": "css"
	};

    this.init = function() {
		_initACE();
        _bindUIEvents();
    };
    
    var _findElements = function() {
        _elements.content = document.getElementById('frm-edit__inp-content');
		_elements.contentType = document.getElementById('frm-edit__inp-contenttype');
        _elements.form = document.getElementById('frm-edit'); 
    };

	var _copyEditorContentToTextarea = function() {
    	_elements.content.value = _elements.editor.getSession().getValue();
	};
    
    var _initEditor = function() {
    	_elements.content.style.display = 'none';
    
    	_elements.editor = ace.edit("frm-edit__cnt-editor");
    	_elements.editor.setTheme("ace/theme/chrome");
    	_elements.editor.getSession().setValue(_elements.content.value);
		_adjustEditorMode();
    	
    	_elements.form.addEventListener('submit', _copyEditorContentToTextarea);
    };

	var _adjustEditorMode = function() {
		for (var contentType in _modesPerContentType) {
			if (_modesPerContentType.hasOwnProperty(contentType) &&
					_contentTypeContains(contentType)) {
				_elements.editor.getSession().setMode('ace/mode/' +
						_modesPerContentType[contentType]);
			}
		}
	};

	var _contentTypeContains = function(s) {
		return _elements.contentType.value.indexOf(s) >= 0;
	};

	var _initUI = function() {
		_findElements();
		_initEditor();
	};

    var _bindUIEvents = function() {
        document.addEventListener('DOMContentLoaded', _initUI);
    };

	var _initACE = function() {
		// TODO: Improve this SOMEHOW
    	ace.config.setModuleUrl("ace/theme/chrome", "/js/ace/theme-chrome.js?template");
    	ace.config.setModuleUrl("ace/mode/javascript", "/js/ace/mode-javascript.js?template");
    	ace.config.setModuleUrl("ace/mode/html", "/js/ace/mode-html.js?template");
    	ace.config.setModuleUrl("ace/mode/css", "/js/ace/mode-css.js?template");
	};
    
})();
GoneEditor.init();

