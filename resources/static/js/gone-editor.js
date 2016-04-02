
var GoneEditor = new(function() {

	var _root = this;
    var _elements = {
        content: null,
        editor: null,
        form: null
    };

    this.init = function() {
        _bindUIEvents();
    };
    
    var _findElements = function() {
        _elements.content = document.getElementById('frm-edit__inp-content');
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
    	
    	_elements.form.addEventListener('submit', _copyEditorContentToTextarea);
    };

	var _initUI = function() {
		_findElements();
		_initEditor();
	};

    var _bindUIEvents = function() {
        document.addEventListener('DOMContentLoaded', _initUI);
    };
    
})();
GoneEditor.init();

