jQuery(document).ready(function() {
  "use strict";
  // Smooth Scroll Scripts
  // $('.menu-scroll li a').mPageScroll2id ({
  //   offset: 120,
  //   scrollSpeed: 1500,
  // });
  // Circle Progress Script
  $('.fmf-circle-progress').circleProgress ({
    startAngle: 300,
    lineCap: 'round',
    animation: {
      duration: 1500,
    },
    emptyFill:'#F3F3F3',
  }).on('circle-animation-progress', function (event, progress, stepValue) {
    $(this).find('.circle-counter').text((stepValue * 100).toFixed(0)+'%');
  });

  // Progress Bar Script
  var delay = 500;
  $('.progress-bar').each(function() {
    $(this).delay(delay).animate ({
      width: $(this).attr('aria-valuenow') + '%'
    }, delay);
    $(this).prop('Counter',0).animate ({
      Counter: $(this).text()
    },
    {
      duration: delay,
      easing: 'swing',
      complete: function() {
        $('.progress-tooltip').fadeIn();
      },
    });
  });

  // Malihu Custom Scrollbar Script
  $('.fmf-scroll').mCustomScrollbar ({
    mouseWheelPixels: 150,
    scrollInertia: 500,
  });

  // Fileuploader Script
  $('.fmf-uploader input[type="file"]').fileuploader ({
    // onSelect: function(item) {
    // if (!item.html.find('.fileuploader-action-start').length)
    //   item.html.find('.fileuploader-action-remove').before('<a class="fileuploader-action fileuploader-action-start" title="Upload"><i></i></a>');
    // },
    upload: {
      // url: 'php/ajax_upload_file.php',
      data: null,
      type: 'POST',
      enctype: 'multipart/form-data',
      start: false,
      synchron: true,
      captions: {
        drop: 'Drag and Drop Files',
      },
      onSuccess: function(result, item) {
        item.html.find('.fileuploader-action-remove').addClass('fileuploader-action-success');
      },
      onError: function(item) {
        item.upload.status != 'cancelled' && item.html.find('.fileuploader-action-retry').length == 0 ? item.html.find('.column-actions').prepend (
          '<a class="fileuploader-action fileuploader-action-retry" title="Retry"><i></i></a>'
        ) : null;
      },
      onComplete: null,
    },
    // onRemove: function(item) {
    //   $.post('./php/ajax_remove_file.php', {
    //     file: item.name
    //   });
    // }
  });

  // Tooltip Script
  $('[data-toggle="tooltip"]').tooltip();

  // Add And Remove Classes Script
  $('.fmf-like').click(function () {
    $('.fmf-like').removeClass('active');
    $(this).addClass('active');   
  });

  // Owl Carousel Slider Script
  $('.owl-carousel').each( function() {
    var $carousel = $(this);
    var $items = ($carousel.data('items') !== undefined) ? $carousel.data('items') : 1;
    var $items_lg = ($carousel.data('items-lg') !== undefined) ? $carousel.data('items-lg') : 1;
    var $items_md = ($carousel.data('items-md') !== undefined) ? $carousel.data('items-md') : 1;
    var $items_sm = ($carousel.data('items-sm') !== undefined) ? $carousel.data('items-sm') : 1;
    var $items_ssm = ($carousel.data('items-ssm') !== undefined) ? $carousel.data('items-ssm') : 1;
    $carousel.owlCarousel ({
      loop : ($carousel.data('loop') !== undefined) ? $carousel.data('loop') : true,
      items : $carousel.data('items'),
      margin : ($carousel.data('margin') !== undefined) ? $carousel.data('margin') : 0,
      dots : ($carousel.data('dots') !== undefined) ? $carousel.data('dots') : true,
      nav : ($carousel.data('nav') !== undefined) ? $carousel.data('nav') : false,
      navText : ["<div class='slider-no-current'><span class='current-no'></span><span class='total-no'></span></div><span class='current-monials'></span>", "<div class='slider-no-next'></div><span class='next-monials'></span>"],
      autoplay : ($carousel.data('autoplay') !== undefined) ? $carousel.data('autoplay') : false,
      autoplayTimeout : ($carousel.data('autoplay-timeout') !== undefined) ? $carousel.data('autoplay-timeout') : 5000,
      animateIn : ($carousel.data('animatein') !== undefined) ? $carousel.data('animatein') : false,
      animateOut : ($carousel.data('animateout') !== undefined) ? $carousel.data('animateout') : false,
      mouseDrag : ($carousel.data('mouse-drag') !== undefined) ? $carousel.data('mouse-drag') : true,
      autoWidth : ($carousel.data('auto-width') !== undefined) ? $carousel.data('auto-width') : false,
      autoHeight : ($carousel.data('auto-height') !== undefined) ? $carousel.data('auto-height') : false,
      center : ($carousel.data('center') !== undefined) ? $carousel.data('center') : false,
      responsiveClass: true,
      dotsEachNumber: true,
      smartSpeed: 600,
      autoplayHoverPause: true,
      responsive : {
        0 : {
          items : $items_ssm,
        },
        480 : {
          items : $items_sm,
        },
        768 : {
          items : $items_md,
        },
        992 : {
          items : $items_lg,
        },
        1200 : {
          items : $items,
        }
      }
    });
    var totLength = $('.owl-dot', $carousel).length;
    $('.total-no', $carousel).html(totLength);
    $('.current-no', $carousel).html(totLength);
    $carousel.owlCarousel();
    $('.current-no', $carousel).html(1);
    $carousel.on('changed.owl.carousel', function(event) {
      var total_items = event.page.count;
      var currentNum = event.page.index + 1;
      $('.total-no', $carousel ).html(total_items);
      $('.current-no', $carousel).html(currentNum);
    });
  });

  // Counter Script
  $('.counter').countUp ({
    delay: 1,
    time: 100,
  });

  // Masonry Script
  var $grid = $('.fmf-masonry').isotope ({
    itemSelector: '.masonry-item',
    layoutMode: 'packery',
    percentPosition: true,
    isFitWidth: true,
  });

  // Number Input Modify Script
  $('.quantity').each(function() {
    var spinner = jQuery(this),
    input = spinner.find('input[type="number"]'),
    btnUp = spinner.find('.quantity-up'),
    btnDown = spinner.find('.quantity-down'),
    min = input.attr('min'),
    max = input.attr('max');
    btnUp.click(function() {
      var oldValue = parseFloat(input.val());
      if (oldValue >= max) {
        var newVal = oldValue;
      }
      else {
        var newVal = oldValue + 1;
      }
      spinner.find('input').val(newVal);
      spinner.find('input').trigger('change');
    });
    btnDown.click(function() {
      var oldValue = parseFloat(input.val());
      if (oldValue <= min) {
        var newVal = oldValue;
      }
      else {
        var newVal = oldValue - 1;
      }
      spinner.find('input').val(newVal);
      spinner.find('input').trigger('change');
    });
  });

  // On Control Change Add And Remove Class In Table Row Script
  $('#selectAll').on('click',function(){
    if(this.checked){
      $(".has-table tbody tr").addClass("active");
      $('input[type=checkbox]').each(function(){
        this.checked = true;
      });
    }else{
      $(".has-table tbody tr").removeClass("active");
      $('input[type=checkbox]').each(function(){
        this.checked = false;
      });
    }
  });
  $('table.has-table td:first-child input[type=checkbox]').click(function () {
    $(this).parents('tr').toggleClass('active');   
  });

  // Toggle Class Script
  $('.aside-toggle').click(function() {
    $(this).parents('.fmf-sidebar').toggleClass('aside-open');
    $('body').toggleClass('manage-page');
  });

  // Outside Click Remove Classe's Script
  $(document).on('click', function(e) {
    if ($(e.target).is('.fmf-aside, .fmf-aside *, .aside-toggle, .aside-toggle *') === false) {
      $('body').removeClass('manage-page');
      $('.fmf-sidebar').removeClass('aside-open');
    }
  });

  // Header Toggle Class Script
  $('.header-aside-toggle').click(function() {
    $(this).parents('.fmf-documentation').toggleClass('aside-main-open');
    $('body').toggleClass('has-sidebar-body');
  });

  // Outside Click Remove Classe's Script
  $(document).on('click', function(e) {
    if ($(e.target).is('.fmf-aside, .fmf-aside *, .header-aside-toggle, .header-aside-toggle *') === false) {
      $('body').removeClass('has-sidebar-body');
      $('.fmf-documentation').removeClass('aside-main-open');
    }
  });

  // Code Box Copy Script
   function copyToClipboard(text, el) {
    var copyTest = document.queryCommandSupported('copy');
    var elOriginalText = el.attr('data-original-title');

    if (copyTest === true) {
      var copyTextArea = document.createElement("textarea");
      copyTextArea.value = text;
      document.body.appendChild(copyTextArea);
      copyTextArea.select();
      try {
        var successful = document.execCommand('copy');
        var msg = successful ? 'Copied!' : 'Whoops, not copied!';
        el.attr('data-original-title', msg).tooltip('show');
      } catch (err) {
        console.log('Oops, unable to copy');
      }
      document.body.removeChild(copyTextArea);
      el.attr('data-original-title', elOriginalText);
    } else {
      // Fallback if browser doesn't support .execCommand('copy')
      //window.prompt("Copy to clipboard: Ctrl+C or Command+C, Enter", text);
    }
  }

  // Tooltip Scripts
  $('.js-tooltip').tooltip();
 
  // Copy To Clipboard Scripts
  $('.js-tooltip').click(function() {
    var text = $(this).parents('.code-box-copy').children('.copy-clipboard').text();
    var el = $(this);
    copyToClipboard(text, el);
  });

  // Datepicker Script
  flatpickr('.fmf-datepicker', {
    locale: {
      weekdays: {
        shorthand: ['S', 'M', 'T', 'W', 'T', 'F', 'S'],
        longhand: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
      },
    },
  });
 // Sticky Scripts
 $(window).resize(function() {
    if ($(window).width() >= 768) {
      $('.fmf-sticky').sticky ({
        topSpacing: 0,
        zIndex: 4
      });
    };
  });
  $(window).trigger('resize');
  
  // Date Rangepicker Script
  flatpickr('.fmf-date-rangepicker', {
    mode: 'range',
    locale: {
      weekdays: {
        shorthand: ['S', 'M', 'T', 'W', 'T', 'F', 'S'],
        longhand: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
      },
    },
  });

  // Multidate Rangepicker Script
  flatpickr('.fmf-multidate-rangepicker', {
    mode: 'range',
    showMonths: 2,
    locale: {
      weekdays: {
        shorthand: ['S', 'M', 'T', 'W', 'T', 'F', 'S'],
        longhand: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
      },
    },
  });
});