package main

import (
	"golang.org/x/net/html/charset"

	"bufio"
	"bytes"
	"encoding/xml"
	"strconv"
	"time"
)

//string | int8 | uint8 | int16 | uint16 | int32 | uint32 | float | double | timestamp
const (
	TAG_CLUSTER        = "CLUSTER"
	TAG_HOST           = "HOST"
	TAG_METRIC         = "METRIC"
	ATTR_NAME          = "NAME"
	ATTR_TYPE          = "TYPE"
	ATTR_VAL           = "VAL"
	TYPE_VAL_STRING    = "string"
	TYPE_VAL_TIMESTAMP = "timestamp"
)

func (f *Fetcher) fetch(metricCh chan []byte, metrics, mb *bytes.Buffer) error {
	r := bufio.NewReader(f.conn)

	ts := time.Now().Unix()
	tsstr := strconv.FormatInt(ts, 10)

	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = charset.NewReaderLabel

	idx := 0
	var namespaces [][]byte
	namespaces = append(namespaces, []byte(f.graphitePrefix))

L:
	for {
		select {
		case <-f.ctx.Done():
			break L
		default:

			t, err := decoder.Token()
			if t == nil {
				break L
			}
			if err != nil {
				Logger.Log("fetch", "xml", "err", err)
				return err
			}

			switch se := t.(type) {
			case xml.StartElement:
				switch se.Name.Local {
				case TAG_CLUSTER:
					for _, attr := range se.Attr {
						if attr.Name.Local == ATTR_NAME {
							namespaces = append(namespaces, []byte(attr.Value))
						}
					}
				case TAG_HOST:
					for _, attr := range se.Attr {
						if attr.Name.Local == ATTR_NAME {
							namespaces = append(namespaces, []byte(attr.Value))
						}
					}
				case TAG_METRIC:
					metric := makeMetric(&se, mb, namespaces, tsstr)
					if len(metric) > 0 {
						metrics.Write(metric)
						metrics.WriteString("\n")
						idx++
					}
				}

			case xml.EndElement:
				if se.Name.Local == TAG_CLUSTER || se.Name.Local == TAG_HOST {
					namespaces = namespaces[:len(namespaces)-1]
				} else if se.Name.Local == TAG_METRIC {
					mb.Reset()
				}
			}

			if idx >= f.flushCnt {
				bs := make([]byte, metrics.Len())
				if _, err := metrics.Read(bs); err != nil {
					Logger.Log("err", err)
				}
				metricCh <- bs
				metrics.Reset()
				idx = 0
			}
		}
	}

	if metrics.Len() > 0 {
		bs := make([]byte, metrics.Len())
		if _, err := metrics.Read(bs); err != nil {
			Logger.Log("err", err)
		}
		metricCh <- bs
		metrics.Reset()
	}

	return nil
}

func makeMetric(el *xml.StartElement, mb *bytes.Buffer, ns [][]byte, ts string) []byte {
	for _, attr := range el.Attr {
		if attr.Name.Local == ATTR_NAME {
			mb.Write(bytes.Join(ns, []byte(".")))
			mb.WriteString(".")
			mb.WriteString(attr.Value)
		} else if attr.Name.Local == ATTR_VAL {
			mb.WriteString(" ")
			mb.WriteString(attr.Value)
			mb.WriteString(" ")
			mb.WriteString(ts)
		} else if attr.Name.Local == ATTR_TYPE {
			if attr.Value != TYPE_VAL_STRING && attr.Value != "" {
				bs := make([]byte, mb.Len())
				if _, err := mb.Read(bs); err != nil {
					Logger.Log("err", err)
				}
				return bs
			}
		}
	}
	return nil
}
