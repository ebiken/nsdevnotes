### encodeVarint(), decodeVarint(), dropMSB() is explained in the BLOG below:
### https://engineering.mercari.com/blog/entry/20210921-ca19c9f371/

### TODO: port this golang code to python.

import sys, getopt

def encodeVarint(val uint64) (int, []byte) {
    length := 1
    for v := val; v >= 1<<7; v >>= 7 {
        length++
    }

    b := make([]byte, 0, length)
    for i := 0; i < length; i++ {
        v := val >> (7 * i) & 0x7f
        if i+1 != length {
            v |= 0x80
        }

        b = append(b, byte(v))
    }

    return len(b), b
}

def decodeVarint(in io.ByteReader) (length int, n uint64, _ error) {
    for i := 0; ; i++ {
        b, err := in.ReadByte()
        if err != nil {
            return 0, 0, err
        }

        length++

        v, hasNext := dropMSB(b)
        n |= uint64(v) << (7 * i)

        if !hasNext {
            return length, n, nil
        }
    }
}

def dropMSB(b byte) (_ byte, hasNext bool) {
    hasNext = b>>7 == 1
    return b & 0x7f, hasNext
}

if __name__ == "__main__":
    result = 0
    try:
        opts, args = getopt.getopt(sys.argv[1:], "e:d:")
    except getopt.GetoptError:
        print("varint.py -e <num> or -d <num>")
        sys.exit(2)
    print(opts)
    for opt, arg in opts:
        if opt == "-e":
            result = encodeVarint(arg)
        elif opt == "-d":
            result = decodeVarint(arg)
    print(result)

